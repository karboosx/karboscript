function main() {
    int a = 1;
    int b = 1;
    int max = readInt();
    while (b < max) {
                out (b);

        int c = b;
        int b = a + b;
        a = c;
    }
}
