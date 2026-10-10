#include <iostream>
#include <vector>

// Helper level 2 function
int calculateHash(int data) {
    return data * 31;
}

// Helper level 1 function
int getBufferSize(int input) {
    int hash = calculateHash(input);
    return hash + 1024;
}

// Core processing logic
void processData() {
    int size = getBufferSize(5);
    std::cout << "Buffer Size: " << size << std::endl;
}

// Entry point
int main() {
    processData();
    return 0;
}