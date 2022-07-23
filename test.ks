function main() {
    $a = 1;
    $b = 1;

    while ($a < 5) {
        $b = 1;
        while ($b < 5) {
            test ($a, $b);
            $b=$b+1;
        }

        $a=$a+1;
    }
}

function test($a, $b) {
    $a = $a+10;
    $b = $b+10;
    out ($a, $b);
}