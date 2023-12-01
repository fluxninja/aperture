package extractors

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	classificationv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
)

// CompileToRego compiles the extractors into a rego.
func CompileToRego(
	packageName string,
	labelExtractors map[string]*classificationv1.Extractor,
) (string, error) {
	if !isValidPackageName(packageName) {
		return "", BadPackageName
	}

	// Sorting for stable output
	sortedLabels := make([]string, 0, len(labelExtractors))
	for label := range labelExtractors {
		sortedLabels = append(sortedLabels, label)
	}
	sort.Strings(sortedLabels)

	out := strings.Builder{}
	fmt.Fprintf(&out, "package %s\n", packageName)
	needs := needs{}
	for _, label := range sortedLabels {
		if !IsRegoIdent(label) {
			return "", fmt.Errorf(
				"%w: %s (allowed chars are alphanumeric and underscore, cannot be Rego keyword)",
				BadLabelName, label,
			)
		}
		compiledExtractor, err := emitExtractor(label, labelExtractors[label], &needs)
		if err != nil {
			return "", fmt.Errorf("%w for label %s: %v", BadExtractor, label, err)
		}
		out.WriteString(compiledExtractor)
	}
	if needs.Segments {
		out.WriteString(segments)
	}
	if needs.Bearer {
		out.WriteString(bearer)
	}
	return out.String(), nil
}

// BadExtractor occurs when extractor is invalid.
var BadExtractor = badExtractor{}

func (b badExtractor) Error() string { return "invalid extractor" }

type badExtractor struct{}

// BadLabelName is an error occurring when label name is invalid.
var BadLabelName = badLabelName{}

type badLabelName struct{}

func (b badLabelName) Error() string { return "invalid label name" }

// BadPackageName occurs when package name is invalid.
var BadPackageName = badPackageName{}

func (b badPackageName) Error() string { return "invalid package name" }

type badPackageName struct{}

type needs struct {
	Segments bool
	Bearer   bool
}

const segments = `
_ninja_components := split(input.attributes.request.http.path, "?")
_ninja_path := _ninja_components[0]
_ninja_segments := split(trim(_ninja_path, "/"), "/")
`

const bearer = `
_ninja_bearer := value {
	header := input.attributes.request.http.headers.authorization
	startswith(header, "Bearer ")
	value := substring(header, count("Bearer "), -1)
}
`

func emitExtractor(key string, e *classificationv1.Extractor, needs *needs) (string, error) {
	if e.GetVariant() == nil {
		return "", errors.New("no variant set")
	}
	switch variant := e.GetVariant().(type) {
	case *classificationv1.Extractor_From:
		from := ParseAttributePath(variant.From)
		if err := from.validate(); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s := %s\n", key, renderAttributePath(from, needs)), nil

	case *classificationv1.Extractor_Json:
		from := ParseAttributePath(variant.Json.From)
		if err := from.validate(); err != nil {
			return "", err
		}
		pointer, err := ParseJSONPointer(variant.Json.Pointer)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"%s := json.unmarshal(%s)%s\n",
			key,
			renderAttributePath(from, needs),
			renderJSONPointer(pointer),
		), nil

	case *classificationv1.Extractor_Jwt:
		from := ParseAttributePath(variant.Jwt.From)
		if err := from.validate(); err != nil {
			return "", err
		}
		pointer, err := ParseJSONPointer(variant.Jwt.JsonPointer)
		if err != nil {
			return "", err
		}
		out := strings.Builder{}
		fmt.Fprintf(&out, "%s := payload%s {\n", key, renderJSONPointer(pointer))
		fmt.Fprintf(
			&out,
			"  [_, payload, _] := io.jwt.decode(%s)\n",
			renderAttributePath(from, needs),
		)
		fmt.Fprintf(&out, "}\n")
		return out.String(), nil

	case *classificationv1.Extractor_Address:
		from := ParseAttributePath(variant.Address.From)
		if err := from.validateAddress(); err != nil {
			return "", err
		}
		out := strings.Builder{}
		fmt.Fprintf(&out, "%s = result {\n", key)
		fmt.Fprintf(&out, "  value := %s\n", renderAttributePath(from, needs))
		fmt.Fprintf(&out, `  result := concat(":", [value.socketAddress.address, format_int(value.socketAddress.portValue, 10)])`)
		fmt.Fprintln(&out)
		fmt.Fprintf(&out, "}\n")
		fmt.Fprintf(&out, "%s = %s.pipe.path\n", key, renderAttributePath(from, needs))
		return out.String(), nil

	case *classificationv1.Extractor_PathTemplates:
		return compilePathTemplates(key, variant.PathTemplates, needs)
	default:
		return "", errors.New("unsupported extractor variant")
	}
}

func renderAttributePath(path AttributePath, needs *needs) string {
	if path.isBearer() {
		needs.Bearer = true
		return "_ninja_bearer"
	} else {
		return "input.attributes" + renderPath([]string(path), attributePathMode)
	}
}

func renderJSONPointer(ptr JSONPointer) string {
	return renderPath(ptr.segments, jsonPathMode)
}

func renderPath(path []string, mode pathMode) string {
	out := strings.Builder{}
	for _, segment := range path {
		if IsRegoIdent(segment) {
			fmt.Fprintf(&out, ".%s", segment)
		} else if mode == jsonPathMode && isUint(segment) {
			fmt.Fprintf(&out, `[{"%s", %s}[_]]`, segment, segment)
		} else {
			fmt.Fprintf(&out, `[%q]`, segment)
		}
	}
	return out.String()
}

type pathMode int

const (
	jsonPathMode pathMode = iota
	attributePathMode
)

func isUint(s string) bool {
	_, err := strconv.ParseUint(s, 10, 64)
	return err == nil
}

func compilePathTemplates(
	key string,
	pts *classificationv1.PathTemplateMatcher,
	needs *needs,
) (string, error) {
	needs.Segments = true

	sortedPTs := make([]ptAndValue, 0, len(pts.TemplateValues))
	for rawPT, value := range pts.TemplateValues {
		pt, err := ParsePathTemplate(rawPT)
		if err != nil {
			return "", err
		}
		sortedPTs = append(sortedPTs, ptAndValue{PTString: pt.String(), PT: pt, Value: value})
	}

	// first sort by template's textual representation, to have stable output
	sort.Sort(byTemplate(sortedPTs))

	// then put the more specific templates first (as order matters)
	sort.Stable(fromMostSpecific(sortedPTs))

	parts := make([]string, 0, len(sortedPTs))
	for _, pt := range sortedPTs {
		parts = append(
			parts,
			fmt.Sprintf("= %q %s", pt.Value, renderPathTemplateRuleBody(pt.PT)),
		)
	}
	return fmt.Sprintf("%s %s\n", key, strings.Join(parts, " else ")), nil
}

type ptAndValue struct {
	PTString string
	Value    string
	PT       PathTemplate
}

type byTemplate []ptAndValue

func (a byTemplate) Len() int      { return len(a) }
func (a byTemplate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byTemplate) Less(i, j int) bool {
	return a[i].PTString < a[j].PTString
}

type fromMostSpecific []ptAndValue

func (a fromMostSpecific) Len() int      { return len(a) }
func (a fromMostSpecific) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a fromMostSpecific) Less(i, j int) bool {
	return a[i].PT.IsMoreSpecificThan(a[j].PT)
}

// Renders a Rego rule body (the `{}` block) that will check if given path template matches.
func renderPathTemplateRuleBody(template PathTemplate) string {
	out := strings.Builder{}
	out.WriteString("{\n")
	lenComparator := "=="
	if template.HasTrailingWildcard {
		lenComparator = ">="
	}
	fmt.Fprintf(&out, "  count(_ninja_segments) %s %d\n", lenComparator, template.NumSegments())
	for i, segment := range template.Prefix {
		fmt.Fprintf(&out, "  _ninja_segments[%d] == %q\n", i, segment)
	}
	out.WriteString("}")
	return out.String()
}
