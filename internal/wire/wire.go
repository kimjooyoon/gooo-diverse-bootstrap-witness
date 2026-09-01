// Package wire contains only the interchange schema shared by the two paths.
// It deliberately contains no parsing, lowering, execution, comparison, or
// policy logic.
package wire

type Meta struct {
	Schema              string              `json:"schema"`
	Authority           string              `json:"authority"`
	ContractID          string              `json:"contract_id"`
	Language            LanguageSpec        `json:"language"`
	Semantic            SemanticSpec        `json:"semantic"`
	GenerationGraph     []GraphNode         `json:"generation_graph"`
	Independence        IndependenceSpec    `json:"independence_predicate"`
	CanonicalComparison CanonicalComparison `json:"canonical_semantic_comparison"`
	TerminalTracePolicy TerminalTracePolicy `json:"terminal_trace_policy"`
	Resolution          ResolutionSpec      `json:"resolution"`
	MeasurementPolicy   MeasurementPolicy   `json:"measurement_policy"`
	FixedCases          []FixedCase         `json:"fixed_cases"`
	OptionalInputs      []OptionalInput     `json:"optional_inputs"`
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
}

type FixedCase struct {
	ID             string         `json:"id"`
	Source         string         `json:"source"`
	ExpectedStatus string         `json:"expected_status"`
	PathBAvailable bool           `json:"path_b_available"`
	PathBVariant   string         `json:"path_b_variant"`
	Unknown        *UnknownRecord `json:"unknown,omitempty"`
	ReplayPathA    bool           `json:"replay_path_a,omitempty"`
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
	Available      bool   `json:"available"`
	SemanticDigest string `json:"semantic_digest,omitempty"`
	ArtifactDigest string `json:"artifact_digest,omitempty"`
	TraceDigest    string `json:"trace_digest,omitempty"`
	ArtifactPath   string `json:"artifact_path,omitempty"`
	IRPath         string `json:"ir_path,omitempty"`
	TracePath      string `json:"trace_path,omitempty"`
}

type IdentityIndicators struct {
	SemanticIdentity      bool  `json:"semantic_identity"`
	ArtifactByteIdentity  bool  `json:"artifact_byte_identity"`
	TerminalTraceIdentity bool  `json:"terminal_trace_identity"`
	RuntimeIdentity       *bool `json:"runtime_identity,omitempty"`
}

type CaseObservation struct {
	ID             string              `json:"id"`
	ExpectedStatus string              `json:"expected_status"`
	ActualStatus   string              `json:"actual_status"`
	SourceDigest   string              `json:"source_digest"`
	PathA          PathObservation     `json:"path_a"`
	PathB          PathObservation     `json:"path_b"`
	ReplayIdentity *IdentityIndicators `json:"replay_identity,omitempty"`
	Identity       IdentityIndicators  `json:"identity"`
	Reason         string              `json:"reason"`
	Unknown        *UnknownRecord      `json:"unknown,omitempty"`
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
	Schema         string            `json:"schema"`
	ContractDigest string            `json:"contract_digest"`
	SuiteStatus    string            `json:"suite_status"`
	FixedCaseCount int               `json:"fixed_case_count"`
	Cases          []CaseObservation `json:"cases"`
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
	WallMS     int64 `json:"wall_ms"`
	PeakRSSKiB int64 `json:"peak_rss_kib"`
}

type RuntimeMetrics struct {
	Build          RuntimeMetric `json:"build"`
	Test           RuntimeMetric `json:"test"`
	Conformance    RuntimeMetric `json:"conformance"`
	GeneratedBuild RuntimeMetric `json:"generated_build"`
	GeneratedRun   RuntimeMetric `json:"generated_run"`
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
	AutoCommitAuthority             int `json:"auto_commit_authority"`
	AutoPushAuthority               int `json:"auto_push_authority"`
	AutoMergeAuthority              int `json:"auto_merge_authority"`
}

type Evidence struct {
	Schema             string                  `json:"schema"`
	ContractDigest     string                  `json:"contract_digest"`
	SuiteStatus        string                  `json:"suite_status"`
	Conformance        ConformanceReport       `json:"conformance"`
	Inventory          InventoryMetrics        `json:"inventory"`
	Runtime            RuntimeMetrics          `json:"runtime"`
	Tests              TestMetrics             `json:"tests"`
	GeneratedArtifacts GeneratedMetrics        `json:"generated_artifacts"`
	RuntimeExecution   RuntimeExecutionMetrics `json:"runtime_execution"`
	Authority          AuthorityMetrics        `json:"authority"`
}
