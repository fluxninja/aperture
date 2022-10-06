package iface

// ClassificationEngine is the interface for registering classifiers.
type ClassificationEngine interface {
	RegisterClassifier(classifier Classifier) error
	UnregisterClassifier(classifier Classifier) error
	GetClassifier(classifierID ClassifierID) (Classifier, error)
}
