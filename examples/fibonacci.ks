function main() {
    $a = 1;
    $b = 1;

    while ($b < 500) {
        
        $c = $b;
        $b = $a + $b;
        $a = $c;
        out ($b);
    }
}