package user

const (
	PRIVILEGED string = "PRIVILEGED"
	BANNED     string = "BANNED"
	PREFIX     string = "PREFIX"
)

type Feature struct {
	Id            string
	FeatureType   string
	FeatureDetail string
	ValidUntil    string
	IsActive      bool
}

type FeatureOption func(Feature) Feature

func WithFeatureId(id string) FeatureOption {
	return func(feature Feature) Feature {
		feature.Id = id
		return feature
	}
}

func WithFeatureType(featureType string) FeatureOption {
	return func(feature Feature) Feature {
		feature.FeatureType = featureType
		return feature
	}
}

func WithFeatureDetail(featureDetail string) FeatureOption {
	return func(feature Feature) Feature {
		feature.FeatureDetail = featureDetail
		return feature
	}
}

func WithValidUntil(validUntil string) FeatureOption {
	return func(feature Feature) Feature {
		feature.ValidUntil = validUntil
		return feature
	}
}

func WithIsActive(isActive bool) FeatureOption {
	return func(feature Feature) Feature {
		feature.IsActive = isActive
		return feature
	}
}

func CreateDefaultFeature() Feature {
	return Feature{}
}

func CreateFeature(options ...FeatureOption) Feature {
	feature := CreateDefaultFeature()

	for _, option := range options {
		feature = option(feature)
	}

	return feature
}
