import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

// Since __dirname is not available in ES modules, we derive it like this:
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Path to your package.json
const packageJsonPath = path.join(__dirname, "package.json");
const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, "utf8"));
const currentVersion = packageJson.version;

// Path to your consts.ts file
const constsPath = path.join(__dirname, "sdk", "consts.ts");
let constsFileContents = fs.readFileSync(constsPath, "utf8");

// Regular expression to match the version string in consts.ts
const versionRegEx = /export const LIBRARY_VERSION = ".*";/;

// Replace the version string with the current version
constsFileContents = constsFileContents.replace(
  versionRegEx,
  `export const LIBRARY_VERSION = "${currentVersion}";`,
);

console.log(constsPath)
console.log(constsFileContents)

// Write the changes to consts.ts
fs.writeFileSync(constsPath, constsFileContents, "utf8");

console.log(`Updated consts.ts with version: ${currentVersion}`);
