package registry

type CartType string

func (c CartType) ToString() string {
	return string(c)
}

const (
	CartTypeDelivery      CartType = "DELIVERY"
	CartTypeUnspecified   CartType = "Unspecified"
	CartTypeDining        CartType = "DINING"
	CartTypeGroupOrdering CartType = "GROUP_ORDERING"
)

type CartLandingSource string

const (
	CartLandingSourceUnspecified      CartLandingSource = "Unspecified"
	CartLandingSourceMenu             CartLandingSource = "Menu"
	CartLandingSourceExpress          CartLandingSource = "Express"
	CartLandingSourceAerobar          CartLandingSource = "Aerobar"
	CartLandingSourceIntercityMenu    CartLandingSource = "IntercityMenu"
	CartLandingSourceIntercityExpress CartLandingSource = "IntercityExpress"
	CartLandingSourceIntercityAerobar CartLandingSource = "IntercityAerobar"
	CartLandingSourceHealthy          CartLandingSource = "Healthy"
)

func (cl CartLandingSource) ToString() string {
	return string(cl)
}

type CartPhase string

func (cp CartPhase) ToString() string {
	return string(cp)
}

const (
	CartPhaseUnspecified CartPhase = "Unspecified"
	CartPhaseInitialise  CartPhase = "Initialise"
	CartPhaseBuild       CartPhase = "Build"
	CartPhaseCheckout    CartPhase = "Checkout"
)

type GroupPhase string

func (gp GroupPhase) ToString() string {
	return string(gp)
}

const (
	GroupPhaseUnspecified          GroupPhase = "Unspecified"
	GroupPhaseCreateGroup          GroupPhase = "CreateGroup"
	GroupPhaseValidateGroup        GroupPhase = "ValidateGroup"
	GroupPhaseGetGroupOrderPreview GroupPhase = "GetGroupOrderPreview"
	GroupPhaseAddMemberToGroup     GroupPhase = "AddMemberToGroup"
)

type ContextMetaData string

const (
	ContextMetaDataAppRequestUUID             ContextMetaData = "app-request-uuid"
	ContextMetaDataVerifyActiveDuplicateOrder ContextMetaData = "verify-active-duplicate-order"
	ContextMetaDataRiderAppInstalled          ContextMetaData = "rider-app-installed"
	ContextMetaDataIsZoman                    ContextMetaData = "is-zoman"
	ContextMetadataShadowRequest              ContextMetaData = "x-shadow-request"
	ContextMetadataShadowRequestId            ContextMetaData = "x-shadow-request-id"
	ContextMetadataTestFlow                   ContextMetaData = "x-test-flow"
	ContextMetadataLandingSource              ContextMetaData = "x-landing-source"
	ContextMetadataBenchmarkingTraffic        ContextMetaData = "benchmarking_traffic"
	ContextMetaDataIsGatewayReqSource         ContextMetaData = "is-gateway-request-source"
)

func (cmd ContextMetaData) ToString() string {
	return string(cmd)
}

type CartSessionTokenKey string

const (
	CartSessionTokenKeyUnspecified          CartSessionTokenKey = "Unspecified"
	CartSessionTokenKeyBenefitsToken        CartSessionTokenKey = "BenefitsToken"
	CartSessionTokenKeyChargesComputationId CartSessionTokenKey = "ChargeComputationId"
)

func (cstk CartSessionTokenKey) ToString() string {
	return string(cstk)
}

func (cstk CartSessionTokenKey) IsEmpty() bool {
	return cstk == CartSessionTokenKeyUnspecified
}

type CartModificationType string

const (
	CartModificationTypeUnspecified       CartModificationType = "Unspecified"
	CartModificationTypeAddItem           CartModificationType = "AddItem"
	CartModificationTypeAddItemMultiStore CartModificationType = "AddItemMultiStore"
	CartModificationTypeRescueOrder       CartModificationType = "RescueOrder"
)

func (cmt CartModificationType) IsEmpty() bool {
	return cmt == CartModificationTypeUnspecified
}

func (cmt CartModificationType) ToString() string {
	return string(cmt)
}

type InvocationMode string

const (
	InvocationModeUnspecified InvocationMode = "Unspecified"
	InvocationModeForeground  InvocationMode = "Foreground"
	InvocationModeBackground  InvocationMode = "Background"
)

func (im InvocationMode) ToString() string {
	return string(im)
}

type InvocationEntityId string

const (
	InvocationEntityIdUnspecified                InvocationEntityId = "Unspecified"
	InvocationEntityIdSimilarCarts               InvocationEntityId = "SimilarCarts"
	InvocationEntityIdOfferWall                  InvocationEntityId = "OfferWall"
	InvocationEntityIdBenefitServiceGetDiscounts InvocationEntityId = "BenefitServiceGetDiscounts"
	InvocationEntityIdReplacementMealSimilarRes  InvocationEntityId = "ReplacementMealSimilarRes"
)

func (iei InvocationEntityId) ToString() string {
	return string(iei)
}

type InvocationReason string

const (
	InvocationReasonUnspecified                        InvocationReason = "Unspecified"
	InvocationReasonGTCalculationPostBenefitSuggestion InvocationReason = "GTCalculationPostBenefitSuggestion"
)

func (ir InvocationReason) ToString() string {
	if ir == InvocationReasonUnspecified {
		return ""
	}

	return string(ir)
}

func NewInvocationReason(reason string) InvocationReason {
	if reason == "" {
		return InvocationReasonUnspecified
	}

	return InvocationReason(reason)
}

type InvocationEntityType string

const (
	InvocationEntityTypeUnspecified InvocationEntityType = "Unspecified"
	InvocationEntityTypeCustomer    InvocationEntityType = "Customer"
	InvocationEntityTypeSystem      InvocationEntityType = "System"
)

func (iet InvocationEntityType) ToString() string {
	return string(iet)
}

type ComputationMode string

const (
	ComputationModeUnspecified                       ComputationMode = "Unspecified"
	ComputationModeFinalValueCalculation             ComputationMode = "FinalValueCalculation"
	ComputationModeUpdateGroupOrderMemberCatalog     ComputationMode = "UpdateGroupOrderMemberCatalog"
	ComputationModeUpdateGroupOrderAggregatedCatalog ComputationMode = "UpdateGroupOrderAggregatedCatalog"
)

func (cm ComputationMode) ToString() string {
	return string(cm)
}

type FeatureSupportKey string

const (
	FeatureSupportKeyUnspecified                   FeatureSupportKey = "Unspecified"
	FeatureSupportKeyZomatoMoneyV2                 FeatureSupportKey = "ZomatoMoneyV2"
	FeatureSupportKeyCartV19                       FeatureSupportKey = "CartV19"
	FeatureSupportKeyUserSalt                      FeatureSupportKey = "UserSalt"
	FeatureSupportKeyBXGYDeduplication             FeatureSupportKey = "BXGYDeduplication"
	FeatureSupportKeyNonPromoItemsBlockedForOffers FeatureSupportKey = "NonPromoItemsBlockedForOffers"
	FeatureSupportKeySaltForSubtotalFix            FeatureSupportKey = "SaltForSubtotalFix"
	FeatureSupportKeyNewBrandReferralFlow          FeatureSupportKey = "NewBrandReferralFlow"
	FeatureSupportKeyStepperOffer                  FeatureSupportKey = "StepperOffer"
	FeatureSupportKeyWebRouteAPI                   FeatureSupportKey = "WebRouteAPI"
	FeatureSupportKeyFreebie                       FeatureSupportKey = "Freebie"
	FeatureSupportKeyNewSaltOffers                 FeatureSupportKey = "NewSaltOffers"
	FeatureSupportKeyBaseDishDiscount              FeatureSupportKey = "BaseDishDiscount"
	FeatureSupportKeyShowRawFeatureOnCart          FeatureSupportKey = "ShowRawFeatureOnCart"
	FeatureSupportKeyDeal                          FeatureSupportKey = "Deal"
	FeatureSupportKeyCartRoundOff                  FeatureSupportKey = "CartRoundOff"
	FeatureSupportKeyAutoApplyZcredits             FeatureSupportKey = "AutoApplyZcredits"
	FeatureSupportKeyCorporateVouchers             FeatureSupportKey = "CorporateVouchers"
	FeatureSupportKeyCombinedAdditiveOffers        FeatureSupportKey = "CombinedAdditiveVouchers"
	FeatureSupportKeyGold                          FeatureSupportKey = "Gold"
	FeatureSupportKeyGoldSubsription               FeatureSupportKey = "GoldSubsription"
	FeatureSupportKeyEnterpriseMealWallet          FeatureSupportKey = "EnterpriseMealWallet"
	FeatureSupportKeyPremiumCheckout               FeatureSupportKey = "PremiumCheckout"
	FeatureSupportKeyGoldLite                      FeatureSupportKey = "GoldLite"
	FeatureSupportKeyBrunch                        FeatureSupportKey = "Brunch"
	FeatureSupportKeyAddOnCustomizationStepper     FeatureSupportKey = "AddOnCustomizationStepper"
	FeatureSupportKeyUnifiedBalance                FeatureSupportKey = "UnifiedBalance"
	FeatureSupportKeyTimedOffer                    FeatureSupportKey = "TimedOffer"
	FeatureSupportKeyEnterpriseLimit               FeatureSupportKey = "EnterpriseLimit"
	FeatureSupportKeyTimedOfferBackgroundMode      FeatureSupportKey = "TimedOfferBackgroundMode"
	FeatureSupportKeyPaymentsBasedPromo            FeatureSupportKey = "PaymentsBasedPromo"
	FeatureSupportKeyTip                           FeatureSupportKey = "Tip"
	FeatureSupportKeyFeedingIndiaV2                FeatureSupportKey = "FeedingIndiaV2"
	FeatureSupportKeyUnifiedZBalancePopupV2        FeatureSupportKey = "UnifiedZBalancePopupV2"
	FeatureSupportKeyPriorityDeliveryV2            FeatureSupportKey = "PriorityDeliveryV2"
	FeatureSupportKeyPayAfterOrder                 FeatureSupportKey = "PayAfterOrder"
	FeatureSupportKeyCarbonOffsetFee               FeatureSupportKey = "CarbonOffsetFee"
	FeatureSupportKeyRiderWelfareFund              FeatureSupportKey = "RiderWelfareFund"
	FeatureSupportKeyVegModeV2                     FeatureSupportKey = "VegCharge"
	FeatureSupportKeyFeedingIndiaRobinhood         FeatureSupportKey = "FeedingIndiaRobinhood"
	FeatureSupportKeyPayAfterOrderV2               FeatureSupportKey = "PayAfterOrderV2"
	FeatureSupportKeyPostOrderPayment              FeatureSupportKey = "PostOrderPayment"
	FeatureSupportKeyMultipleGroupOrderSupport     FeatureSupportKey = "MultipleGroupOrderSupport"
	FeatureSupportKeyDynamicDietaryTag             FeatureSupportKey = "DynamicDietaryTag"
)

func (k FeatureSupportKey) ToString() string {
	return string(k)
}

type EntryPointType string

const (
	EntryPointTypeUnspecified            EntryPointType = "unspecified"
	EntryPointTypeSearchLargeOrderBanner EntryPointType = "searchLargeOrderBanner"
)

func (e EntryPointType) ToString() string {
	if e == EntryPointTypeUnspecified {
		return ""
	}

	return string(e)
}

func (e EntryPointType) FromString(val string) EntryPointType {
	if val == "" {
		return EntryPointTypeUnspecified
	}

	return EntryPointType(val)
}
