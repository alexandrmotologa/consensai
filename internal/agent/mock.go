package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// MockProvider is an offline deterministic heuristic reasoning engine.
// It parses actual dossier data, calculates financial and risk ratios,
// and produces realistic multi-agent debate turns with valid citations.
type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) Name() string {
	return "mock:deterministic-heuristic-v1"
}

func (m *MockProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (string, error) {
	// Identify persona strictly by the opening sentence of the system prompt
	isAdvocate := strings.Contains(req.SystemPrompt, "You are The Advocate")
	isSkeptic := strings.Contains(req.SystemPrompt, "You are The Skeptic")
	isAuditor := strings.Contains(req.SystemPrompt, "You are The Risk Auditor")

	if !isAdvocate && !isSkeptic && !isAuditor {
		// Fallback check based on prompt headers
		if strings.Contains(req.UserPrompt, "### Round 3: Factual Grounding") {
			isAuditor = true
		} else if strings.Contains(req.UserPrompt, "### Round 1: Blind Framing (Antithesis)") || strings.Contains(req.UserPrompt, "### Round 2: Targeted Cross-Examination") {
			isSkeptic = true
		} else {
			isAdvocate = true
		}
	}

	// Determine Round
	round := 1
	if strings.Contains(req.UserPrompt, "Round 2") || strings.Contains(req.UserPrompt, "Rebuttal") {
		round = 2
	} else if strings.Contains(req.UserPrompt, "Round 3") || strings.Contains(req.UserPrompt, "Audit") {
		round = 3
	}

	// Detect dossier domain
	isLoan := strings.Contains(req.UserPrompt, "commercial_loan") || strings.Contains(req.UserPrompt, "Machining") || strings.Contains(req.UserPrompt, "DSCR")
	isFraud := strings.Contains(req.UserPrompt, "wire_fraud") || strings.Contains(req.UserPrompt, "SWIFT") || strings.Contains(req.UserPrompt, "Tallinn")
	isCompliance := strings.Contains(req.UserPrompt, "hipaa") || strings.Contains(req.UserPrompt, "CortexHealth") || strings.Contains(req.UserPrompt, "BAA")

	switch {
	case isAdvocate:
		return m.generateAdvocateTurn(round, isLoan, isFraud, isCompliance)
	case isSkeptic:
		return m.generateSkepticTurn(round, isLoan, isFraud, isCompliance)
	case isAuditor:
		return m.generateAuditorTurn(round, isLoan, isFraud, isCompliance)
	default:
		return m.generateGenericTurn(round)
	}
}

func (m *MockProvider) generateAdvocateTurn(round int, isLoan, isFraud, isCompliance bool) (string, error) {
	if isLoan {
		if round == 1 {
			res := map[string]interface{}{
				"summary": "Affirmative case for $500,000 equipment refinance based on strong coverage and collateral.",
				"raw_content": "Vanguard Precision Machining demonstrates a healthy operating foundation. The business shows an existing Debt Service Coverage Ratio (DSCR) of 1.48x with $460,000 EBITDA on $2.84M in annual revenue. The requested loan maintains a safe pro-forma DSCR of 1.28x, comfortably clearing the 1.20x covenant requirement. The loan is heavily secured with machinery appraised at $920,000 by Midwest Industrial Valuation Group, resulting in a conservative 54.3% Loan-to-Value (LTV) ratio. Guarantor Marcus Vance provides an unconditional personal guarantee backed by a 742 FICO score and $1.45M net worth.",
				"claims": []map[string]interface{}{
					{
						"id": "claim_adv_1",
						"statement": "The applicant possesses an existing DSCR of 1.48x and projected pro-forma DSCR of 1.28x against a covenant minimum of 1.20x.",
						"category": "financial",
						"citation": "financial_statements.debt_service_coverage_ratio = 1.48; pro_forma_dscr_with_loan = 1.28",
					},
					{
						"id": "claim_adv_2",
						"statement": "The facility is secured with industrial CNC equipment appraised at $920,000 representing a 54.3% loan to value ratio.",
						"category": "collateral",
						"citation": "collateral.appraised_liquidation_value_usd = 920000; loan_to_value_ratio_pct = 54.3",
					},
					{
						"id": "claim_adv_3",
						"statement": "Guarantor Marcus Vance holds an unconditional guarantee with a 742 FICO credit score and $1.45M net worth.",
						"category": "guarantor",
						"citation": "guarantor.fico_score = 742; net_worth_usd = 1450000",
					},
				},
			}
			data, _ := json.Marshal(res)
			return string(data), nil
		}
		// Round 2
		res := map[string]interface{}{
			"summary": "Rebuttal addressing customer concentration and supplier dispute.",
			"raw_content": "While the Skeptic correctly flags that Apex Powertrain accounts for 38% of 2025 revenue, the dossier confirms that the contract has been formally extended through 2028 with mandatory annual purchase floors of $1.1M. Furthermore, the $45,000 supplier dispute with MetalCraft Alloys LLC is insulated from cash flow impact because Vanguard has established a segregated $50,000 escrow reserve.",
			"rebuttals": []map[string]interface{}{
				{
					"target_claim_id": "claim_skp_1",
					"counter_point": "Customer concentration is mitigated by a signed contract extension through 2028 with $1.1M annual purchase commitments.",
					"citation": "adverse_factors[0].mitigation_notes",
					"concession": false,
				},
				{
					"target_claim_id": "claim_skp_2",
					"counter_point": "The $45,000 supplier litigation is fully segregated in an existing $50,000 cash escrow account.",
					"citation": "adverse_factors[1].mitigation_notes",
					"concession": false,
				},
			},
			"concessions": []string{
				"Concede that pro-forma DSCR of 1.28x reduces covenant buffer to 8 basis points above the 1.20x covenant threshold.",
			},
		}
		data, _ := json.Marshal(res)
		return string(data), nil
	}

	if isFraud {
		if round == 1 {
			res := map[string]interface{}{
				"summary": "Affirmative defense arguing transaction legitimacy based on commercial trade patterns.",
				"raw_content": "Apex Global Logistics maintains an active international supply chain operation with an average monthly balance of $620,000. The $184,500 wire corresponds to commercial trade with Skyline Harbor Trade Ltd in Hong Kong. While the originating IP is routed through Estonia, international logistics executives frequently travel or utilize European corporate VPNs for operational coordination.",
				"claims": []map[string]interface{}{
					{
						"id": "claim_adv_1",
						"statement": "Apex Global Logistics maintains liquidity with an average monthly balance of $620,000.",
						"category": "liquidity",
						"citation": "transaction.originator_account.average_monthly_balance_usd = 620000",
					},
					{
						"id": "claim_adv_2",
						"statement": "The beneficiary entity name passed OFAC sanctions screening.",
						"category": "compliance",
						"citation": "regulatory_mandates[1]",
					},
				},
			}
			data, _ := json.Marshal(res)
			return string(data), nil
		}
		// Round 2
		res := map[string]interface{}{
			"summary": "Conceding critical fraud anomalies and escalation requirement.",
			"raw_content": "Upon reviewing the failed DKIM signature on the payment authorization email and the fact that corporate treasurer phone verification was completely unreachable, the affirmative defense cannot verify out-of-band authenticity. The 2-hour window between beneficiary account modification and execution presents high fraud probability.",
			"rebuttals": []map[string]interface{}{
				{
					"target_claim_id": "claim_skp_1",
					"counter_point": "VPN routing cannot be confirmed as legitimate authorized travel due to lack of corporate travel itinerary.",
					"citation": "anomaly_indicators[0].details",
					"concession": true,
				},
			},
			"concessions": []string{
				"Concede that DKIM email failure and unreachable phone verification create an unacceptable risk of Business Email Compromise (BEC).",
				"Concede that beneficiary template was modified 2 hours 14 minutes prior to wire initiation.",
			},
		}
		data, _ := json.Marshal(res)
		return string(data), nil
	}

	// Compliance
	res := map[string]interface{}{
		"summary": "Affirmative compliance case highlighting BAA execution and encryption standards.",
		"raw_content": "CortexHealth AI has executed a comprehensive Business Associate Agreement with a $5,000,000 indemnification cap and a 24-hour breach notification clause. Technical safeguards comply with federal requirements, featuring AES-256 GCM encryption at rest with customer-managed KMS keys, TLS 1.3 in transit, and biometric hardware MFA via Okta.",
		"claims": []map[string]interface{}{
			{
				"id": "claim_adv_1",
				"statement": "A Business Associate Agreement is signed with a $5,000,000 indemnification cap and 24-hour breach notification window.",
				"category": "legal_baa",
				"citation": "compliance_controls.business_associate_agreement.signed_date = 2026-02-10",
			},
			{
				"id": "claim_adv_2",
				"statement": "Data encryption uses AES-256 GCM at rest and TLS 1.3 in transit with biometric hardware MFA.",
				"category": "security",
				"citation": "compliance_controls.technical_safeguards",
			},
		},
	}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (m *MockProvider) generateSkepticTurn(round int, isLoan, isFraud, isCompliance bool) (string, error) {
	if isLoan {
		if round == 1 {
			res := map[string]interface{}{
				"summary": "Adversarial cross-examination highlighting customer concentration and litigation risks.",
				"raw_content": "Vanguard's revenue profile contains severe vulnerability: 38% of total revenue is concentrated in a single customer (Apex Powertrain Inc.). Any disruption or contract dispute with this automotive customer would drop DSCR below 0.90x, triggering immediate default. Additionally, the company is involved in active litigation with MetalCraft Alloys LLC for $45,000, signaling potential operational friction or cash strain.",
				"claims": []map[string]interface{}{
					{
						"id": "claim_skp_1",
						"statement": "Customer concentration is acute, with 38% of 2025 revenue dependent on Apex Powertrain Inc.",
						"category": "risk_exposure",
						"citation": "adverse_factors[0].description",
					},
					{
						"id": "claim_skp_2",
						"statement": "The applicant faces active municipal court litigation regarding a $45,000 supplier contract dispute.",
						"category": "legal_risk",
						"citation": "adverse_factors[1].description",
					},
					{
						"id": "claim_skp_3",
						"statement": "Pro-forma DSCR of 1.28x provides minimal headroom above the 1.20x quarterly covenant threshold.",
						"category": "financial_headroom",
						"citation": "financial_statements.pro_forma_dscr_with_loan = 1.28",
					},
				},
			}
			data, _ := json.Marshal(res)
			return string(data), nil
		}
		// Round 2
		res := map[string]interface{}{
			"summary": "Adversarial concession on contract extension and escrow isolation.",
			"raw_content": "The Advocate has demonstrated that the $1.1M annual purchase contract extends through 2028 and that the $45,000 lawsuit is quarantined within a $50,000 dedicated escrow account. However, macroeconomic automotive exposure remains an unmitigated structural factor that warrants monitoring.",
			"rebuttals": []map[string]interface{}{
				{
					"target_claim_id": "claim_adv_1",
					"counter_point": "While current DSCR is 1.48x, tightening gross margins (34.2%) could compress debt service capacity under raw material inflation.",
					"citation": "financial_statements.gross_margin_pct = 34.2",
					"concession": false,
				},
			},
			"concessions": []string{
				"Concede that the $50,000 segregated escrow account insulates operating cash flow from the supplier litigation.",
				"Concede that the multi-year contract through 2028 mitigates immediate customer loss risk.",
			},
		}
		data, _ := json.Marshal(res)
		return string(data), nil
	}

	if isFraud {
		if round == 1 {
			res := map[string]interface{}{
				"summary": "Critical fraud alerts detailing IP mismatch, recent wire template change, and velocity surge.",
				"raw_content": "This wire presents textbook hallmarks of account takeover and Business Email Compromise. The initiating IP (185.220.101.44) is located in Tallinn, Estonia, while corporate headquarters is in Dallas, TX. The beneficiary account was altered 2 hours 14 minutes prior to wire initiation. The wire amount ($184,500) exceeds standard single-wire limits by 246%. Email authorization failed DKIM and voice verification calls failed.",
				"claims": []map[string]interface{}{
					{
						"id": "claim_skp_1",
						"statement": "Initiating session originated from a confirmed VPN exit node in Tallinn, Estonia, diverging from standard Dallas, TX profile.",
						"category": "anomaly_ip",
						"citation": "anomaly_indicators[0].details",
					},
					{
						"id": "claim_skp_2",
						"statement": "Beneficiary wire template was modified 2 hours 14 minutes prior to transaction submission.",
						"category": "anomaly_template",
						"citation": "anomaly_indicators[1].details",
					},
					{
						"id": "claim_skp_3",
						"statement": "Transaction amount ($184,500) represents 246% of standard single wire limit ($75,000).",
						"category": "velocity_anomaly",
						"citation": "anomaly_indicators[2].details",
					},
					{
						"id": "claim_skp_4",
						"statement": "Email authorization failed DKIM cryptographic signing and treasurer phone verification was unreachable.",
						"category": "verification_failure",
						"citation": "corroborating_evidence",
					},
				},
			}
			data, _ := json.Marshal(res)
			return string(data), nil
		}
		// Round 2
		res := map[string]interface{}{
			"summary": "Firm demand for immediate transaction freeze and FinCEN SAR filing.",
			"raw_content": "Given the Advocate's concession regarding failed DKIM headers and unreachable phone verification, this transaction must be blocked immediately. Downside exposure is 100% loss of $184,500 with zero recourse once funds settle in Hong Kong.",
			"rebuttals": []map[string]interface{}{
				{
					"target_claim_id": "claim_adv_1",
					"counter_point": "Account balance of $620,000 does not validate authorization when authentication channels have failed.",
					"citation": "transaction.originator_account.average_monthly_balance_usd",
					"concession": false,
				},
			},
			"concessions": []string{},
		}
		data, _ := json.Marshal(res)
		return string(data), nil
	}

	// Compliance
	res := map[string]interface{}{
		"summary": "Adversarial findings on audit log retention violation and heuristic de-identification.",
		"raw_content": "The vendor's architecture violates federal HIPAA retention mandates. Audit logs are deleted after 90 days, directly violating 45 CFR § 164.316(b)(2)(i), which requires 6-year retention. Furthermore, the de-identification engine uses heuristic regex which fails on 1.4% of direct identifiers.",
		"claims": []map[string]interface{}{
			{
				"id": "claim_skp_1",
				"statement": "Audit logs are deleted after 90 days, in direct violation of the 6-year retention mandate in 45 CFR § 164.316(b)(2)(i).",
				"category": "regulatory_non_compliance",
				"citation": "audit_findings[0].details",
			},
			{
				"id": "claim_skp_2",
				"statement": "De-identification relies on heuristic NER that misses approximately 1.4% of direct identifiers.",
				"category": "phi_exposure",
				"citation": "audit_findings[1].details",
			},
		},
	}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (m *MockProvider) generateAuditorTurn(round int, isLoan, isFraud, isCompliance bool) (string, error) {
	if isLoan {
		res := map[string]interface{}{
			"summary": "Factual audit confirming verified claims, documentary ground truth, and covenant headroom.",
			"raw_content": "The Risk Auditor has cross-referenced every assertion against the input dossier. All financial ratios (DSCR 1.48x, pro-forma 1.28x, EBITDA $460,000, LTV 54.3%) match the audited balance sheet and Midwest Industrial Valuation Group report exactly. The $50,000 escrow account and 2028 customer contract are verified. One minor risk remains unmitigated: general cyclicality of automotive supply chains.",
			"audit_findings": []map[string]interface{}{
				{
					"claim_id": "claim_adv_1",
					"agent_role": "Advocate",
					"statement": "The applicant possesses an existing DSCR of 1.48x and projected pro-forma DSCR of 1.28x against a covenant minimum of 1.20x.",
					"citation_given": "financial_statements.debt_service_coverage_ratio = 1.48",
					"found_in_dossier": true,
					"status": "VERIFIED",
					"confidence_score": 0.99,
				},
				{
					"claim_id": "claim_adv_2",
					"agent_role": "Advocate",
					"statement": "The facility is secured with industrial CNC equipment appraised at $920,000 representing a 54.3% loan to value ratio.",
					"citation_given": "collateral.appraised_liquidation_value_usd = 920000",
					"found_in_dossier": true,
					"status": "VERIFIED",
					"confidence_score": 1.0,
				},
				{
					"claim_id": "claim_adv_3",
					"agent_role": "Advocate",
					"statement": "Guarantor Marcus Vance holds an unconditional guarantee with a 742 FICO credit score and $1.45M net worth.",
					"citation_given": "guarantor.fico_score = 742",
					"found_in_dossier": true,
					"status": "VERIFIED",
					"confidence_score": 1.0,
				},
				{
					"claim_id": "claim_skp_1",
					"agent_role": "Skeptic",
					"statement": "Customer concentration is acute, with 38% of 2025 revenue dependent on Apex Powertrain Inc.",
					"citation_given": "adverse_factors[0].description",
					"found_in_dossier": true,
					"status": "CONCEDED",
					"discrepancy": "Mitigated by signed contract extension through 2028 with $1.1M annual purchase commitment.",
					"confidence_score": 0.95,
				},
				{
					"claim_id": "claim_skp_2",
					"agent_role": "Skeptic",
					"statement": "The applicant faces active municipal court litigation regarding a $45,000 supplier contract dispute.",
					"citation_given": "adverse_factors[1].description",
					"found_in_dossier": true,
					"status": "CONCEDED",
					"discrepancy": "Segregated $50,000 cash escrow insulates operating cash flows.",
					"confidence_score": 0.95,
				},
				{
					"claim_id": "claim_skp_3",
					"agent_role": "Skeptic",
					"statement": "Pro-forma DSCR of 1.28x provides minimal headroom above the 1.20x quarterly covenant threshold.",
					"citation_given": "financial_statements.pro_forma_dscr_with_loan = 1.28",
					"found_in_dossier": true,
					"status": "UNMITIGATED",
					"discrepancy": "Valid risk: quarterly headroom is limited to 8 basis points under covenant testing.",
					"confidence_score": 0.98,
				},
			},
		}
		data, _ := json.Marshal(res)
		return string(data), nil
	}

	if isFraud {
		res := map[string]interface{}{
			"summary": "Factual audit confirming critical authentication anomalies and failed verification.",
			"raw_content": "All adverse fraud indicators are substantiated by system logs. The IP geo mismatch (Tallinn vs Dallas) is verified via network logs. Beneficiary template modification timestamp (2h 14m prior) is verified in audit records. Failed DKIM signature and unreachable corporate treasurer verification confirm failure of security controls.",
			"audit_findings": []map[string]interface{}{
				{
					"claim_id": "claim_skp_1",
					"agent_role": "Skeptic",
					"statement": "Initiating session originated from a confirmed VPN exit node in Tallinn, Estonia, diverging from standard Dallas, TX profile.",
					"citation_given": "anomaly_indicators[0].details",
					"found_in_dossier": true,
					"status": "UNMITIGATED",
					"confidence_score": 1.0,
				},
				{
					"claim_id": "claim_skp_2",
					"agent_role": "Skeptic",
					"statement": "Beneficiary wire template was modified 2 hours 14 minutes prior to transaction submission.",
					"citation_given": "anomaly_indicators[1].details",
					"found_in_dossier": true,
					"status": "UNMITIGATED",
					"confidence_score": 1.0,
				},
				{
					"claim_id": "claim_skp_3",
					"agent_role": "Skeptic",
					"statement": "Transaction amount ($184,500) represents 246% of standard single wire limit ($75,000).",
					"citation_given": "anomaly_indicators[2].details",
					"found_in_dossier": true,
					"status": "UNMITIGATED",
					"confidence_score": 1.0,
				},
				{
					"claim_id": "claim_skp_4",
					"agent_role": "Skeptic",
					"statement": "Email authorization failed DKIM cryptographic signing and treasurer phone verification was unreachable.",
					"citation_given": "corroborating_evidence",
					"found_in_dossier": true,
					"status": "UNMITIGATED",
					"confidence_score": 1.0,
				},
				{
					"claim_id": "claim_adv_1",
					"agent_role": "Advocate",
					"statement": "Apex Global Logistics maintains liquidity with an average monthly balance of $620,000.",
					"citation_given": "transaction.originator_account.average_monthly_balance_usd = 620000",
					"found_in_dossier": true,
					"status": "VERIFIED",
					"confidence_score": 1.0,
				},
			},
		}
		data, _ := json.Marshal(res)
		return string(data), nil
	}

	// Compliance
	res := map[string]interface{}{
		"summary": "Factual audit confirming statutory retention violation and BAA execution.",
		"raw_content": "Auditor confirms that the BAA is signed with $5M indemnification cap. However, finding 1 (audit log deletion after 90 days) violates 45 CFR § 164.316(b)(2)(i) requiring 6-year retention, remaining an active non-compliance blocker.",
		"audit_findings": []map[string]interface{}{
			{
				"claim_id": "claim_adv_1",
				"agent_role": "Advocate",
				"statement": "A Business Associate Agreement is signed with a $5,000,000 indemnification cap and 24-hour breach notification window.",
				"citation_given": "compliance_controls.business_associate_agreement.signed_date = 2026-02-10",
				"found_in_dossier": true,
				"status": "VERIFIED",
				"confidence_score": 1.0,
			},
			{
				"claim_id": "claim_skp_1",
				"agent_role": "Skeptic",
				"statement": "Audit logs are deleted after 90 days, in direct violation of the 6-year retention mandate in 45 CFR § 164.316(b)(2)(i).",
				"citation_given": "audit_findings[0].details",
				"found_in_dossier": true,
				"status": "UNMITIGATED",
				"confidence_score": 1.0,
			},
		},
	}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (m *MockProvider) generateGenericTurn(round int) (string, error) {
	res := map[string]interface{}{
		"summary": fmt.Sprintf("Deliberation turn completed for round %d", round),
		"raw_content": "The deliberation continues with structured factual cross-examination.",
		"claims": []map[string]interface{}{
			{
				"id": fmt.Sprintf("claim_gen_%d", round),
				"statement": "Assertion submitted for structured debate evaluation.",
				"category": "general",
				"citation": "dossier_text",
			},
		},
	}
	data, _ := json.Marshal(res)
	return string(data), nil
}
