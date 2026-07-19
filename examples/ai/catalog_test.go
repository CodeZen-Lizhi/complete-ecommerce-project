package ai

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var exerciseDirectoryPattern = regexp.MustCompile(`^\d{2}_[a-z0-9_]+$`)
var todoPattern = regexp.MustCompile(`TODO ([0-9]+)`)

var requiredExerciseSymbols = map[string][]string{
	"phase3/01_embedding_similarity":           {"embedder", "parseEmbeddingInput", "loadEmbeddingConfig", "cosineSimilarity", "rankTopK"},
	"phase3/02_document_loading_metadata":      {"documentLoader", "pdfLoaderConfig", "newPDFPageLoader", "validateDocumentPath", "normalizeDocument"},
	"phase3/03_chunking_strategies":            {"chunker", "buildParentChildChunks"},
	"phase3/04_vector_index_and_retrieval":     {"indexer", "retriever", "qdrantConfig", "newQdrantHTTPClient", "validateIndexBatch"},
	"phase3/05_hnsw_tuning":                    {"hnswIndex", "evaluateHNSW"},
	"phase3/06_parent_child_retrieval":         {"childRetriever", "parentStore", "expandParents"},
	"phase3/07_hybrid_search_rrf":              {"denseRetriever", "sparseRetriever", "hybridBackendConfig", "newHybridRetrievers", "reciprocalRankFusion"},
	"phase3/08_query_rewrite_and_filter":       {"queryRewriter", "filteredRetriever", "validateQueryFilter"},
	"phase3/09_rerank":                         {"reranker", "rerankerConfig", "newRemoteReranker", "validateRerankResults"},
	"phase3/10_context_budget_and_citations":   {"tokenCounter", "selectContext"},
	"phase3/11_index_lifecycle":                {"documentStateStore", "vectorIndex", "newMySQLStateDB", "runImport"},
	"phase3/12_rag_evaluation":                 {"recallAtK", "reciprocalRank", "summarizeRAGEvaluation"},
	"phase4/01_local_readonly_tool":            {"readonlyTool", "toolExecutor", "newEinoToolsNode", "validateToolInfo"},
	"phase4/02_tool_argument_validation":       {"decodeInventoryQuery", "validateInventoryQuery", "validationError"},
	"phase4/03_context_identity_authorization": {"withIdentity", "identityFromContext", "authorizeResource"},
	"phase4/04_ecommerce_query_tools":          {"productQueryService", "orderQueryService", "queryOrderTool"},
	"phase4/05_multi_tool_orchestration":       {"toolRegistry", "modelTurn", "bindRealTools", "runToolLoop"},
	"phase4/06_tool_resilience":                {"classifyToolFailure", "invokeWithRetry"},
	"phase4/07_safe_write_tool":                {"confirmationStore", "idempotencyStore", "executeWriteTool"},
	"phase5/01_chain_workflow":                 {"chainNode", "buildChain"},
	"phase5/02_graph_routing":                  {"classifyRoute", "buildRoutingGraph"},
	"phase5/03_react_agent":                    {"reactModel", "reactToolsNode", "newEinoReActAgent", "runReAct"},
	"phase5/04_agent_state_and_budget":         {"stateStore", "checkBudget", "advanceState"},
	"phase5/05_loop_detection":                 {"normalizeArguments", "recordCall"},
	"phase5/06_interrupt_checkpoint_resume":    {"checkpointStore", "compileInterruptibleGraph", "resumeFromCheckpoint"},
	"phase5/07_async_agent_task":               {"taskRepository", "runWorker"},
	"phase5/08_task_recovery_idempotency":      {"recoveryStore", "idempotencyRecordStore", "recoverTasks"},
	"phase6/01_golden_dataset":                 {"loadGoldenDataset", "validateGoldenCase"},
	"phase6/02_retrieval_metrics":              {"recallAtK", "reciprocalRank", "summarizeRetrieval"},
	"phase6/03_answer_quality_checks":          {"judgeResult", "newJudgeModel", "validateCitations", "evaluateFacts"},
	"phase6/04_tool_and_agent_metrics":         {"evaluateExecution", "executionMetrics"},
	"phase6/05_latency_token_cost_metrics":     {"percentileDuration", "summarizeMeasurements"},
	"phase6/06_regression_comparison":          {"compareSnapshots", "regressionReport"},
	"phase6/07_eval_ci_gate":                   {"evaluateGate", "exitCodeForGate"},
	"phase7/01_secret_and_log_redaction":       {"loadSecretConfig", "redactFields", "containsSecret"},
	"phase7/02_prompt_injection_defense":       {"evaluateInjection", "authorizeToolCall"},
	"phase7/03_rate_limit_and_concurrency":     {"rateLimiter", "concurrencyLimiter", "runLimited"},
	"phase7/04_end_to_end_tracing":             {"tracer", "span", "newOpenTelemetryProvider", "tracePipeline"},
	"phase7/05_metrics_and_alerts":             {"metricRecorder", "validateMetricLabels", "evaluateAlert"},
	"phase7/06_docker_compose_stack":           {"validateComposeFile", "composeCommands"},
	"phase7/07_health_gray_release_rollback":   {"releaseControlConfig", "newReleaseControlClient", "validateHealth", "evaluateRollout"},
	"phase7/08_dependency_failure_fallback":    {"decideFallback", "recoveryController"},
}

// TestExerciseCatalog 验证七个阶段的目录数量、连续编号和必需骨架文件。
func TestExerciseCatalog(t *testing.T) {
	t.Helper()

	expectedCounts := map[int]int{
		1: 7,
		2: 12,
		3: 12,
		4: 7,
		5: 8,
		6: 7,
		7: 8,
	}

	for phase := 1; phase <= 7; phase++ {
		phase := phase
		t.Run(fmt.Sprintf("phase%d", phase), func(t *testing.T) {
			validatePhase(t, phase, expectedCounts[phase])
		})
	}
}

// validatePhase 验证单个阶段的目录顺序，以及 README 和 Go 骨架是否齐全。
func validatePhase(t *testing.T, phase int, expectedCount int) {
	t.Helper()

	phaseDir := fmt.Sprintf("phase%d", phase)
	entries, err := os.ReadDir(phaseDir)
	if err != nil {
		t.Fatalf("读取 %s 失败: %v", phaseDir, err)
	}

	directories := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if !exerciseDirectoryPattern.MatchString(entry.Name()) {
			t.Fatalf("%s 包含不符合 NN_slug 规范的目录 %q", phaseDir, entry.Name())
		}
		directories = append(directories, entry.Name())
	}
	sort.Strings(directories)

	if len(directories) != expectedCount {
		t.Fatalf("%s 目录数量为 %d，期望 %d: %v", phaseDir, len(directories), expectedCount, directories)
	}

	for index, directory := range directories {
		expectedPrefix := fmt.Sprintf("%02d_", index+1)
		if directory[:3] != expectedPrefix {
			t.Errorf("%s 第 %d 个目录为 %q，期望前缀 %q", phaseDir, index+1, directory, expectedPrefix)
		}

		assertRegularFile(t, filepath.Join(phaseDir, directory, "README.md"))
		if phase >= 2 {
			mainPath := filepath.Join(phaseDir, directory, "main.go")
			assertRegularFile(t, mainPath)

			todoPath := mainPath
			if phase >= 3 {
				exercisePath := filepath.Join(phaseDir, directory, "exercise.go")
				assertRegularFile(t, exercisePath)
				todoPath = exercisePath
			}
			validateTODOSequence(t, todoPath)

			key := filepath.Join(phaseDir, directory)
			if symbols, exists := requiredExerciseSymbols[key]; exists {
				validateExerciseSymbols(t, key, symbols)
			}
		}
	}
}

// validateExerciseSymbols 验证练习包包含题目专属接口、类型或函数签名。
func validateExerciseSymbols(t *testing.T, directory string, required []string) {
	t.Helper()

	packages, err := parser.ParseDir(token.NewFileSet(), directory, func(info os.FileInfo) bool {
		return filepath.Ext(info.Name()) == ".go" && !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Errorf("解析练习目录 %s 失败: %v", directory, err)
		return
	}

	found := make(map[string]struct{})
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, declaration := range file.Decls {
				switch value := declaration.(type) {
				case *ast.FuncDecl:
					found[value.Name.Name] = struct{}{}
				case *ast.GenDecl:
					for _, specification := range value.Specs {
						if typeSpec, ok := specification.(*ast.TypeSpec); ok {
							found[typeSpec.Name.Name] = struct{}{}
						}
					}
				}
			}
		}
	}

	for _, symbol := range required {
		if _, exists := found[symbol]; !exists {
			t.Errorf("练习目录 %s 缺少题目专属符号 %q", directory, symbol)
		}
	}
}

// validateTODOSequence 验证 Go 骨架中首次出现的 TODO 编号从 1 连续递增。
func validateTODOSequence(t *testing.T, path string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("读取 TODO 文件 %s 失败: %v", path, err)
		return
	}

	matches := todoPattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		t.Errorf("Go 骨架 %s 没有编号 TODO", path)
		return
	}

	seen := make(map[int]struct{}, len(matches))
	numbers := make([]int, 0, len(matches))
	for _, match := range matches {
		number, err := strconv.Atoi(string(match[1]))
		if err != nil {
			t.Errorf("解析 %s 的 TODO 编号失败: %v", path, err)
			return
		}
		if _, exists := seen[number]; exists {
			continue
		}
		seen[number] = struct{}{}
		numbers = append(numbers, number)
	}

	for index, number := range numbers {
		expected := index + 1
		if number != expected {
			t.Errorf("%s 的首次 TODO 顺序为 %v，期望在位置 %d 出现 TODO %d", path, numbers, index+1, expected)
			return
		}
	}
}

// assertRegularFile 验证指定路径存在且是普通文件。
func assertRegularFile(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("必需文件 %s 不可用: %v", path, err)
		return
	}
	if !info.Mode().IsRegular() {
		t.Errorf("必需路径 %s 不是普通文件", path)
	}
}
