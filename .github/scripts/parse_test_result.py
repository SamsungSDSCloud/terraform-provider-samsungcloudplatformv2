import json
import os
import sys
from datetime import datetime, timezone

def parse_results(json_file):
    tests = {}   # key: (package, test_name)
    outputs = {} # key: (package, test_name) -> list of output lines

    with open(json_file, "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            try:
                entry = json.loads(line)
            except json.JSONDecodeError:
                continue

            action  = entry.get("Action", "")
            pkg     = entry.get("Package", "")
            test    = entry.get("Test", "")
            elapsed = entry.get("Elapsed", 0.0)

            if not test:
                continue

            key = (pkg, test)

            if action in ("pass", "fail", "skip"):
                tests[key] = {"action": action, "elapsed": elapsed, "package": pkg, "test": test}
            elif action == "output":
                outputs.setdefault(key, []).append(entry.get("Output", ""))

    return tests, outputs


def pkg_short(pkg):
    return pkg.split("/")[-1]


def fmt_elapsed(seconds):
    m = int(seconds) // 60
    s = int(seconds) % 60
    return f"{m}분 {s:02d}초"


def icon(action):
    if action == "pass":
        return "PASS"
    if action == "fail":
        return "FAIL"
    return "SKIP"


def build_body(tests, region, run_url, run_date):
    rows = sorted(tests.values(), key=lambda r: (r["package"], r["test"]))

    total  = len(rows)
    passed = sum(1 for r in rows if r["action"] == "pass")
    failed = sum(1 for r in rows if r["action"] == "fail")

    # column widths
    W_SVC  = max((len(pkg_short(r["package"])) for r in rows), default=10)
    W_SVC  = max(W_SVC, 14)
    W_TEST = max((len(r["test"]) for r in rows), default=20)
    W_TEST = max(W_TEST, 20)

    header = (
        f"{'No':>4}  "
        f"{'서비스 (패키지)':<{W_SVC}}  "
        f"{'테스트명':<{W_TEST}}  "
        f"{'결과':<6}  "
        f"소요시간"
    )
    sep = "─" * (4 + 2 + W_SVC + 2 + W_TEST + 2 + 6 + 2 + 8)

    lines = []
    lines.append("Terraform Acceptance Test 결과 보고")
    lines.append(f"실행일시 : {run_date}")
    lines.append(f"리전     : {region}")
    lines.append(f"Run URL  : {run_url}")
    lines.append("=" * 70)
    lines.append("")
    lines.append(header)
    lines.append(sep)

    for i, r in enumerate(rows, 1):
        result_icon = "✅" if r["action"] == "pass" else ("❌" if r["action"] == "fail" else "⏭️")
        lines.append(
            f"{i:>4}  "
            f"{pkg_short(r['package']):<{W_SVC}}  "
            f"{r['test']:<{W_TEST}}  "
            f"{result_icon} {icon(r['action']):<4}  "
            f"{fmt_elapsed(r['elapsed'])}"
        )

    lines.append("")
    lines.append("=" * 70)
    lines.append(f"요약: 총 {total}건 / 성공 {passed} / 실패 {failed}")
    lines.append("=" * 70)

    if failed > 0:
        lines.append("※ 실패 항목의 상세 로그는 첨부파일(test_failures.log)을 확인하세요.")

    return "\n".join(lines)


def build_failure_log(tests, outputs):
    failed = [r for r in tests.values() if r["action"] == "fail"]
    if not failed:
        return None

    lines = []
    for r in sorted(failed, key=lambda r: (r["package"], r["test"])):
        key = (r["package"], r["test"])
        lines.append("=" * 70)
        lines.append(f"패키지 : {r['package']}")
        lines.append(f"테스트 : {r['test']}")
        lines.append("-" * 70)
        for out in outputs.get(key, []):
            lines.append(out.rstrip())
        lines.append("")

    return "\n".join(lines)


def main():
    json_file   = os.environ.get("TEST_RESULT_FILE", "test_result.json")
    body_file   = os.environ.get("MAIL_BODY_FILE",   "email.txt")
    fail_file   = os.environ.get("MAIL_FAILURE_FILE","test_failures.log")
    region      = os.environ.get("SCP_TF_DEFAULT_REGION", "unknown")
    run_url     = os.environ.get("GITHUB_RUN_URL",   "")
    run_date    = os.environ.get("RUN_DATE", datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M UTC"))

    if not os.path.isfile(json_file):
        print(f"ERROR: {json_file} not found", file=sys.stderr)
        sys.exit(1)

    tests, outputs = parse_results(json_file)

    body = build_body(tests, region, run_url, run_date)
    with open(body_file, "w", encoding="utf-8") as f:
        f.write(body)
    print(f"Email body written to {body_file}")

    failure_log = build_failure_log(tests, outputs)
    if failure_log:
        with open(fail_file, "w", encoding="utf-8") as f:
            f.write(failure_log)
        print(f"Failure log written to {fail_file}")
        print(f"attachment={fail_file}")
    else:
        print(f"attachment=")


if __name__ == "__main__":
    main()
