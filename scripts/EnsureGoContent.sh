#!/usr/bin/env bash
set -u

TARGET_PATH="${1:-src}"

if [[ ! -d "$TARGET_PATH" ]]; then
    echo "Target path does not exist or is not a directory: $TARGET_PATH" >&2
    exit 2
fi

if [[ -t 1 ]]; then
    RED=$'\033[31m'
    DARK_GRAY=$'\033[90m'
    CYAN=$'\033[36m'
    GREEN=$'\033[32m'
    RESET=$'\033[0m'
else
    RED=""
    DARK_GRAY=""
    CYAN=""
    GREEN=""
    RESET=""
fi

# These assume standard gofmt formatting where top-level definitions start at the beginning of a line.
pType='^[[:space:]]*type[[:space:]]+(\(|[[:alpha:]_][[:alnum:]_]*)'
pMethod='^[[:space:]]*func[[:space:]]+\('
pHelper='^[[:space:]]*func[[:space:]]+[[:alpha:]_][[:alnum:]_]*'
pConst='^[[:space:]]*const[[:space:]]+'
pVar='^[[:space:]]*var[[:space:]]+'

violations=0

run_check() {
    local file_glob="$1"
    local forbidden="$2"
    local rule_name="$3"

    local file
    local display_file

    while IFS= read -r -d '' file; do
        display_file="${file#"$PWD"/}"
        display_file="${display_file#./}"

        if ! awk \
            -v forbidden="$forbidden" \
            -v ruleName="$rule_name" \
            -v displayFile="$display_file" \
            -v red="$RED" \
            -v darkGray="$DARK_GRAY" \
            -v reset="$RESET" \
'
BEGIN {
    braceDepth = 0
    parenDepth = 0
    inBlockComment = 0
    inRawString = 0
    fileViolations = 0
    sq = sprintf("%c", 39)
}

function count_char(s, target,    i, c, n) {
    c = 0
    for (i = 1; i <= length(s); i++) {
        if (substr(s, i, 1) == target) {
            c++
        }
    }
    return c
}

function strip_go_line(line,    out, i, n, ch, next, endPos, delim, c) {
    out = ""
    i = 1
    n = length(line)

    while (i <= n) {
        if (inBlockComment) {
            endPos = index(substr(line, i), "*/")
            if (endPos == 0) {
                return out
            }

            inBlockComment = 0
            i += endPos + 1
            continue
        }

        if (inRawString) {
            endPos = index(substr(line, i), "`")
            if (endPos == 0) {
                return out
            }

            inRawString = 0
            i += endPos
            continue
        }

        ch = substr(line, i, 1)
        next = ""

        if (i < n) {
            next = substr(line, i + 1, 1)
        }

        # Line comment
        if (ch == "/" && next == "/") {
            break
        }

        # Block comment start
        if (ch == "/" && next == "*") {
            inBlockComment = 1
            i += 2
            continue
        }

        # Raw string start
        if (ch == "`") {
            inRawString = 1
            i++
            continue
        }

        # Interpreted string or rune literal
        if (ch == "\"" || ch == sq) {
            delim = ch
            i++

            while (i <= n) {
                c = substr(line, i, 1)

                if (c == "\\") {
                    i += 2
                    continue
                }

                if (c == delim) {
                    i++
                    break
                }

                i++
            }

            continue
        }

        out = out ch
        i++
    }

    return out
}

{
    originalLine = $0
    stripped = strip_go_line($0)

    if (stripped ~ /^[[:space:]]*$/) {
        next
    }

    # Only check top-level declarations.
    if (braceDepth == 0 && parenDepth == 0) {
        if (stripped ~ forbidden) {
            printf "%s[%s Violation] %s:%d%s\n", red, ruleName, displayFile, FNR, reset
            printf "%s  Found forbidden content: %s%s\n", darkGray, originalLine, reset
            fileViolations = 1
        }
    }

    braceDepth += count_char(stripped, "{") - count_char(stripped, "}")
    parenDepth += count_char(stripped, "(") - count_char(stripped, ")")

    if (braceDepth < 0) {
        braceDepth = 0
    }

    if (parenDepth < 0) {
        parenDepth = 0
    }
}

END {
    exit fileViolations ? 1 : 0
}
' "$file"; then
            violations=1
        fi
    done < <(find "$TARGET_PATH" -type f -name "$file_glob" -print0)
}

printf "%sStarting Go Content Enforcement...%s\n" "$CYAN" "$RESET"

# Rule 1: *types*.go files should not contain methods, helpers, vars, or consts.
run_check "*types*.go" "$pMethod|$pHelper|$pVar|$pConst" "*types*.go"

# Rule 2: *helpers*.go files should not contain types, methods, vars, or consts.
run_check "*helpers*.go" "$pType|$pMethod|$pVar|$pConst" "*helpers*.go"

# Rule 3: *methods*.go files should not contain types, helpers, vars, or consts.
run_check "*methods*.go" "$pType|$pHelper|$pVar|$pConst" "*methods*.go"

# Rule 4: *constants*.go files should not contain types, methods, helpers, or vars.
run_check "*constants*.go" "$pType|$pMethod|$pHelper|$pVar" "*constants*.go"

# Rule 5: *vars*.go files should not contain types, methods, helpers, or consts.
run_check "*vars*.go" "$pType|$pMethod|$pHelper|$pConst" "*vars*.go"

if [[ "$violations" -eq 0 ]]; then
    printf "%sCheck complete.%s\n" "$GREEN" "$RESET"
    exit 0
else
    echo "Check complete with violations."
    exit 1
fi
