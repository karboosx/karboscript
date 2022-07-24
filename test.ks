function main() {
    out("Enter name: ");
    name = test();
    out("Your name is:", name);
}

function test() {
    return readLine();
}