cd "$(dirname "$0")"

VISUAL=0
if [ "$1" = "--visual" ]; then
    VISUAL=1
fi

total_start=$(date +%s.%N)

for day in day{1..12}; do
    if [ -d "$day" ]; then
        echo "========== Running $day =========="
        start=$(date +%s.%N)
        (cd "$day" && AOC_VISUAL=$VISUAL go run .)
        end=$(date +%s.%N)
        elapsed=$(echo "$end - $start" | bc)
        printf "⭐️⭐️  %s completed in %.3f seconds\n\n" "$day" "$elapsed"
    fi
done

total_end=$(date +%s.%N)
total_elapsed=$(echo "$total_end - $total_start" | bc)
printf "========================================\n"
printf "✅ Total time: %.3f seconds\n" "$total_elapsed"