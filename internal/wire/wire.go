// Package wire contains only the interchange schema shared by the two paths.
// It deliberately contains no parsing, lowering, execution, comparison, or
// policy logic.
package wire

type Meta struct {
	Schema              string              `json:"schema"`
	Authority           string              `json:"authority"`
	ContractID          string              `json:"contract_id"`
	ContractEvolution   ContractEvolution   `json:"contract_evolution"`
	Language            LanguageSpec        `json:"language"`
	Semantic            SemanticSpec        `json:"semantic"`
	SemanticKernel      SemanticKernel      `json:"semantic_kernel"`
	CaseEvidencePolicy  CaseEvidencePolicy  `json:"case_evidence_policy"`
	GenerationGraph     []GraphNode         `json:"generation_graph"`
	Independence        IndependenceSpec    `json:"independence_predicate"`
	WitnessPlan         WitnessPlan         `json:"witness_plan"`
	CanonicalComparison CanonicalComparison `json:"canonical_semantic_comparison"`
	TerminalTracePolicy TerminalTracePolicy `json:"terminal_trace_policy"`
	Resolution          ResolutionSpec      `json:"resolution"`
	MeasurementPolicy   MeasurementPolicy   `json:"measurement_policy"`
	FixedCases          []FixedCase         `json:"fixed_cases"`
	OptionalInputs      []OptionalInput     `json:"optional_inputs"`
}

type ContractEvolution struct {
	PreviousSchema      string   `json:"previous_schema"`
	PreviousContractID  string   `json:"previous_contract_id"`
	PreviousDenominator int      `json:"previous_denominator"`
	CurrentDenominator  int      `json:"current_denominator"`
	AppendOnly          bool     `json:"append_only"`
	AddedFixedCaseIDs   []string `json:"added_fixed_case_ids"`
	RetiredFixedCaseIDs []string `json:"retired_fixed_case_ids"`
}

type LanguageSpec struct {
	Name            string            `json:"name"`
	SourceExtension string            `json:"source_extension"`
	CommentPrefixes []string          `json:"comment_prefixes"`
	Keywords        map[string]string `json:"keywords"`
	Grammar         []string          `json:"grammar"`
}

type SemanticSpec struct {
	IRSchema               string   `json:"ir_schema"`
	Fields                 []string `json:"fields"`
	Normalization          []string `json:"normalization"`
	Comparison             string   `json:"comparison"`
	CommentsAreNonsemantic bool     `json:"comments_are_nonsemantic"`
	SourceAuthority        string   `json:"source_authority"`
	ObservationFields      []string `json:"observation_fields"`
}

type SemanticKernel struct {
	Source              KernelSource        `json:"source"`
	CanonicalInput      CanonicalInput      `json:"canonical_input"`
	ExpectedObservation ExpectedObservation `json:"expected_observation"`
	Evaluation          EvaluationSpec      `json:"evaluation"`
}

type KernelSource struct {
	Path        string `json:"path"`
	Digest      string `json:"digest"`
	DigestScope string `json:"digest_scope"`
}

type CanonicalInput struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	Digest      string `json:"digest"`
	DigestScope string `json:"digest_scope"`
}

type ExpectedObservation struct {
	Status             string   `json:"status"`
	SemanticIdentity   bool     `json:"semantic_identity"`
	DecisionIdentity   bool     `json:"decision_identity"`
	ProvenanceIdentity bool     `json:"provenance_identity"`
	RequiredDigests    []string `json:"required_digests"`
}

type EvaluationSpec struct {
	Authority  string         `json:"authority"`
	Operations []string       `json:"operations"`
	Mutations  []MutationSpec `json:"mutations"`
}

type MutationSpec struct {
	ID     string `json:"id"`
	Effect string `json:"effect"`
	Target string `json:"target"`
	Value  string `json:"value"`
}

type GraphNode struct {
	ID             string   `json:"id"`
	Kind           string   `json:"kind"`
	Implementation string   `json:"implementation,omitempty"`
	Input          any      `json:"input"`
	Outputs        []string `json:"outputs"`
}

type IndependenceSpec struct {
	Paths                            []string `json:"paths"`
	AllowedSharedPackages            []string `json:"allowed_shared_packages"`
	AllowedSharedInputs              []string `json:"allowed_shared_inputs"`
	ForbiddenSharedPackageKinds      []string `json:"forbidden_shared_package_kinds"`
	RequireCIImportIntersectionCheck bool     `json:"require_ci_import_intersection_check"`
	RequireDistinctSources           bool     `json:"require_distinct_parser_lowerer_emitter_executor_sources"`
	RequireIndependentWitnesses      bool     `json:"require_independent_witnesses"`
	SameLineageStatus                string   `json:"same_lineage_status"`
	DuplicateArtifactStatus          string   `json:"duplicate_artifact_status"`
	IdentityFields                   []string `json:"identity_fields"`
}

type WitnessPlan struct {
	RequiredIndependentCount int              `json:"required_independent_count"`
	Stage0Reference          WitnessIdentity  `json:"stage0_reference_witness"`
	GeneratedCurrent         WitnessIdentity  `json:"generated_current_witness"`
	OptionalDiverse          *WitnessIdentity `json:"optional_diverse_witness,omitempty"`
	ProofClaims              []ProofClaim     `json:"proof_claims"`
	Trilemma                 TrilemmaNote     `json:"munchhausen_trilemma"`
}

type WitnessIdentity struct {
	ID             string `json:"id"`
	Role           string `json:"role"`
	Implementation string `json:"implementation"`
	SourcePath     string `json:"source_path"`
	SourceDigest   string `json:"source_digest"`
	LineageID      string `json:"lineage_id"`
	ToolchainID    string `json:"toolchain_id"`
	InputIdentity  string `json:"input_identity"`
	Available      bool   `json:"available"`
	RequiredGate   int    `json:"required_gate"`
}

type ProofClaim struct {
	ID                string   `json:"id"`
	Claim             string   `json:"claim"`
	ProofChoice       string   `json:"proof_choice"`
	IndependenceBasis string   `json:"independence_basis"`
	RequiredEvidence  []string `json:"required_evidence"`
}

type CaseEvidencePolicy struct {
	ProofChoices     []string `json:"proof_choices"`
	IndicatorClasses []string `json:"indicator_classes"`
	AuthorityFields  []string `json:"authority_fields"`
	IRFields         []string `json:"ir_fields"`
	VerifierFields   []string `json:"verifier_fields"`
	ReleaseFields    []string `json:"release_fields"`
}

type TrilemmaNote struct {
	Acknowledged bool     `json:"acknowledged"`
	Foundation  string   `json:"foundation"`
	Coherence   string   `json:"coherence"`
	Regression  string   `json:"regression"`
	Limitations []string `json:"limitations"`
}

type CanonicalComparison struct {
	IdentityIndicator        string `json:"identity_indicator"`
	Digest                   string `json:"digest"`
	FirstDifference          string `json:"first_difference"`
	ArtifactCannotSubstitute bool   `json:"artifact_identity_cannot_substitute"`
}

type TerminalTracePolicy struct {
	Schema               string   `json:"schema"`
	CompareFields        []string `json:"compare_fields"`
	EffectOrder          string   `json:"effect_order"`
	MissingTraceStatus   string   `json:"missing_trace_status"`
	DriftStatus          string   `json:"drift_status"`
	ReasonDriftIsFailure bool     `json:"terminal_reason_drift_is_semantic_failure"`
}

type ResolutionSpec struct {
	StatusPrecedence []string `json:"status_precedence"`
	ClosedRequires   []string `json:"closed_requires"`
	UnknownFields    []string `json:"unknown_fields"`
	RefutedReasons   []string `json:"refuted_reasons"`
}

type MeasurementPolicy struct {
	InventoryExcludePaths []string `json:"inventory_exclude_paths"`
	InventoryFields       []string `json:"inventory_fields"`
	RuntimeFields         []string `json:"runtime_fields"`
	TestFields            []string `json:"test_fields"`
	GeneratedFields       []string `json:"generated_fields"`
	ImprovementRule       string   `json:"improvement_rule"`
	ScoreAggregation      string   `json:"score_aggregation"`
	EstimatedRate         string   `json:"estimated_rate"`
	WitnessRuntimeFields  []string `json:"witness_runtime_fields"`
	ImprovementMissingStatus string `json:"improvement_missing_status"`
	ImprovementPairFields []string `json:"improvement_pair_fields"`
}

type FixedCase struct {
	ID                string         `json:"id"`
	Source            string         `json:"source"`
	ExpectedStatus    string         `json:"expected_status"`
	ScenarioClass     string         `json:"scenario_class"`
	ProofChoice       string         `json:"proof_choice"`
	IndicatorClass    string         `json:"indicator_class"`
	PathBAvailable    bool           `json:"path_b_available"`
	PathBVariant      string         `json:"path_b_variant"`
	IdentityMode      string         `json:"identity_mode"`
	DigestClaim       string         `json:"path_b_digest_claim,omitempty"`
	SelfApprovalCycle bool           `json:"self_approval_cycle,omitempty"`
	FrozenBootstrapMismatch bool     `json:"frozen_bootstrap_mismatch,omitempty"`
	FrozenBootstrapDigest string      `json:"frozen_bootstrap_digest,omitempty"`
	Unknown           *UnknownRecord `json:"unknown,omitempty"`
	ReplayPathA       bool           `json:"replay_path_a,omitempty"`
}

type OptionalInput struct {
	Name         string `json:"name"`
	Tag          string `json:"tag"`
	Commit       string `json:"commit"`
	ReleaseAsset string `json:"release_asset"`
	Digest       string `json:"digest"`
	RequiredGate int    `json:"required_gate"`
	Use          string `json:"use"`
}

type Binding struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Diagnostic struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SemanticIR struct {
	Schema      string       `json:"schema"`
	Program     string       `json:"program"`
	CaseID      string       `json:"case_id,omitempty"`
	ProofChoice string       `json:"proof_choice,omitempty"`
	IndicatorClass string    `json:"indicator_class,omitempty"`
	Bindings    []Binding    `json:"bindings"`
	Emissions   []string     `json:"emissions"`
	Effects     []Binding    `json:"effects"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type EffectEvent struct {
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TerminalTrace struct {
	Schema         string        `json:"schema"`
	TerminalReason string        `json:"terminal_reason"`
	Effects        []EffectEvent `json:"effects"`
}

type GeneratedResult struct {
	IR            SemanticIR
	ArtifactBytes []byte
	Trace         TerminalTrace
}

type PathObservation struct {
	Available         bool   `json:"available"`
	SemanticDigest    string `json:"semantic_digest,omitempty"`
	DecisionDigest    string `json:"decision_digest,omitempty"`
	ProvenanceDigest  string `json:"provenance_digest,omitempty"`
	ObservationDigest string `json:"observation_digest,omitempty"`
	ArtifactDigest    string `json:"artifact_digest,omitempty"`
	TraceDigest       string `json:"trace_digest,omitempty"`
	ArtifactPath      string `json:"artifact_path,omitempty"`
	IRPath            string `json:"ir_path,omitempty"`
	TracePath         string `json:"trace_path,omitempty"`
}

type IdentityIndicators struct {
	SemanticIdentity      bool  `json:"semantic_identity"`
	DecisionIdentity      bool  `json:"decision_identity"`
	ProvenanceIdentity    bool  `json:"provenance_identity"`
	ArtifactByteIdentity  bool  `json:"artifact_byte_identity"`
	TerminalTraceIdentity bool  `json:"terminal_trace_identity"`
	RuntimeIdentity       *bool `json:"runtime_identity,omitempty"`
}

type WitnessObservation struct {
	ID                string `json:"id"`
	Role              string `json:"role"`
	Implementation    string `json:"implementation"`
	SourceDigest      string `json:"source_digest"`
	LineageID         string `json:"lineage_id"`
	ToolchainID       string `json:"toolchain_id"`
	InputIdentity     string `json:"input_identity"`
	InputDigest       string `json:"input_digest"`
	SemanticDigest    string `json:"semantic_digest"`
	DecisionDigest    string `json:"decision_digest"`
	ProvenanceDigest  string `json:"provenance_digest"`
	ObservationDigest string `json:"observation_digest"`
	Independent       bool   `json:"independent"`
}

type ProofReceipt struct {
	ClaimID           string `json:"claim_id"`
	Claim             string `json:"claim"`
	ProofChoice       string `json:"proof_choice"`
	IndicatorClass    string `json:"indicator_class,omitempty"`
	Status            string `json:"status"`
	IndependenceBasis string `json:"independence_basis"`
	EvidenceDigest    string `json:"evidence_digest"`
}

type BootstrapObservation struct {
	RequiredIndependentCount int              `json:"required_independent_count"`
	KernelSourceDigest      string           `json:"kernel_source_digest"`
	CanonicalInputDigest    string           `json:"canonical_input_digest"`
	Stage0Reference         WitnessIdentity  `json:"stage0_reference_witness"`
	GeneratedCurrent        WitnessIdentity  `json:"generated_current_witness"`
	OptionalDiverse         *WitnessIdentity `json:"optional_diverse_witness,omitempty"`
	Trilemma                TrilemmaNote     `json:"munchhausen_trilemma"`
}

type Denominator struct {
	ID                string         `json:"id"`
	Version           string         `json:"version"`
	FixedCaseCount    int            `json:"fixed_case_count"`
	PreviousCaseCount int            `json:"previous_case_count"`
	AddedCaseCount    int            `json:"added_case_count"`
	RetiredCaseCount  int            `json:"retired_case_count"`
	StatusCounts      map[string]int `json:"status_counts"`
	CaseIDs           []string       `json:"case_ids"`
	AppendOnly        bool           `json:"append_only"`
}

type CaseObservation struct {
	ID             string               `json:"id"`
	ExpectedStatus string               `json:"expected_status"`
	ActualStatus   string               `json:"actual_status"`
	SourceDigest   string               `json:"source_digest"`
	ProofChoice    string               `json:"proof_choice"`
	IndicatorClass string               `json:"indicator_class"`
	PathA          PathObservation      `json:"path_a"`
	PathB          PathObservation      `json:"path_b"`
	Witnesses      []WitnessObservation `json:"witnesses,omitempty"`
	ReplayIdentity *IdentityIndicators  `json:"replay_identity,omitempty"`
	Identity       IdentityIndicators   `json:"identity"`
	ProofReceipts  []ProofReceipt       `json:"proof_receipts,omitempty"`
	Reason         string               `json:"reason"`
	Unknown        *UnknownRecord       `json:"unknown,omitempty"`
}

type UnknownRecord struct {
	Stage         string `json:"stage"`
	Step          string `json:"step"`
	Reason        string `json:"reason"`
	UnknownClass  string `json:"unknown_class"`
	NextOperation string `json:"next_operation"`
	BlockedBy     string `json:"blocked_by"`
}

type ConformanceReport struct {
	Schema               string            `json:"schema"`
	ContractVersion      string            `json:"contract_version"`
	ContractDigest       string            `json:"contract_digest"`
	KernelSourceDigest   string            `json:"kernel_source_digest"`
	CanonicalInputDigest string            `json:"canonical_input_digest"`
	SuiteStatus          string            `json:"suite_status"`
	FixedCaseCount       int               `json:"fixed_case_count"`
	Denominator          Denominator       `json:"denominator"`
	Bootstrap            BootstrapObservation `json:"bootstrap"`
	ProofReceipts        []ProofReceipt    `json:"proof_receipts"`
	Cases                []CaseObservation `json:"cases"`
}

type InventoryMetrics struct {
	GoFiles           int `json:"go_files"`
	GoooFiles         int `json:"gooo_files"`
	GoPhysicalLines   int `json:"go_physical_lines"`
	GoooPhysicalLines int `json:"gooo_physical_lines"`
	Directories       int `json:"directories"`
	RegularFiles      int `json:"regular_files"`
}

type RuntimeMetric struct {
	WallMS     *int64         `json:"wall_ms"`
	PeakRSSKiB *int64         `json:"peak_rss_kib"`
	Status     string         `json:"status"`
	Unknown    *UnknownRecord `json:"unknown,omitempty"`
}

type RuntimeMetrics struct {
	Build          RuntimeMetric `json:"build"`
	Test           RuntimeMetric `json:"test"`
	Conformance    RuntimeMetric `json:"conformance"`
	GeneratedBuild RuntimeMetric `json:"generated_build"`
	GeneratedRun   RuntimeMetric `json:"generated_run"`
	WallMS         *int64         `json:"wall_ms"`
	PeakRSSKiB     *int64         `json:"peak_rss_kib"`
	BuildMS        *int64         `json:"build_ms"`
	TestMS         *int64         `json:"test_ms"`
	CacheHits      *int64         `json:"cache_hits"`
	CacheMisses    *int64         `json:"cache_misses"`
	CacheStatus    string         `json:"cache_status"`
	CacheUnknown   *UnknownRecord `json:"cache_unknown,omitempty"`
	Status         string         `json:"status"`
	Unknown        *UnknownRecord `json:"unknown,omitempty"`
}

type TestMetrics struct {
	Total    int `json:"total"`
	Selected int `json:"selected"`
	Executed int `json:"executed"`
	Reused   int `json:"reused"`
	Failed   int `json:"failed"`
	Unknown  int `json:"unknown"`
}

type GeneratedMetrics struct {
	Files int   `json:"files"`
	Bytes int64 `json:"bytes"`
}

type RuntimeExecutionMetrics struct {
	AvailableArtifacts int `json:"available_artifacts"`
	ExecutedArtifacts  int `json:"executed_artifacts"`
	FailedArtifacts    int `json:"failed_artifacts"`
	UnknownArtifacts   int `json:"unknown_artifacts"`
}

type AuthorityMetrics struct {
	RepositoryWrites                int `json:"repository_writes"`
	InputRepositoryWrites           int `json:"input_repository_writes"`
	GeneratedOutputInsideRepository int `json:"generated_output_inside_repository"`
	CrossProjectRequiredGates       int `json:"cross_project_required_gates"`
	AutoCommitAuthority             int `json:"auto_commit_authority"`
	AutoPushAuthority               int `json:"auto_push_authority"`
	AutoMergeAuthority              int `json:"auto_merge_authority"`
}

type WitnessRuntimeObservation struct {
	WitnessID   string `json:"witness_id"`
	Role        string `json:"role"`
	ScenarioID  string `json:"scenario_id"`
	InputDigest string `json:"input_digest"`
	ToolchainID string `json:"toolchain_id"`
	WallMS      *int64         `json:"wall_ms"`
	PeakRSSKiB  *int64         `json:"peak_rss_kib"`
	BuildMS     *int64         `json:"build_ms"`
	TestMS      *int64         `json:"test_ms"`
	CacheHits   *int64         `json:"cache_hits"`
	CacheMisses *int64         `json:"cache_misses"`
	CacheStatus string         `json:"cache_status"`
	CacheUnknown *UnknownRecord `json:"cache_unknown,omitempty"`
	Status      string         `json:"status"`
	Unknown     *UnknownRecord `json:"unknown,omitempty"`
}

type ImprovementObservation struct {
	Status      string        `json:"status"`
	Improvement *int64        `json:"improvement"`
	MatchedPair bool          `json:"matched_pair"`
	Unknown     *UnknownRecord `json:"unknown,omitempty"`
}

type Evidence struct {
	Schema             string                     `json:"schema"`
	ContractVersion    string                     `json:"contract_version"`
	ContractDigest     string                     `json:"contract_digest"`
	SuiteStatus        string                     `json:"suite_status"`
	Conformance        ConformanceReport          `json:"conformance"`
	Bootstrap          BootstrapObservation       `json:"bootstrap"`
	Inventory          InventoryMetrics           `json:"inventory"`
	Runtime            RuntimeMetrics             `json:"runtime"`
	WitnessRuntime     []WitnessRuntimeObservation `json:"witness_runtime"`
	Tests              TestMetrics                `json:"tests"`
	GeneratedArtifacts GeneratedMetrics           `json:"generated_artifacts"`
	RuntimeExecution   RuntimeExecutionMetrics    `json:"runtime_execution"`
	Improvement        ImprovementObservation     `json:"improvement"`
	Authority          AuthorityMetrics           `json:"authority"`
}
