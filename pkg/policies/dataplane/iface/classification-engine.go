package iface

//go:generate mockgen -source=classification-engine.go -destination=../../mocks/mock_classification_engine.go -package=mocks

// ClassificationEngine is the interface for registering classifiers.
type ClassificationEngine interface {
	RegisterClassifier(classifier Classifier) error
	UnregisterClassifier(classifier Classifier) error
	GetClassifier(classifierID ClassifierID) (Classifier, error)
}
