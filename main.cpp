#include<iostream>
#include "engine.h"
#include "utils.h"
using namespace std;

    void f4(){
        cout << "hello4\n";
    }
    void f5(){
        cout << "hello5\n";
    }

    void f2 (){
        cout << "Hello2\n";
    }

    void f1(){
        cout << "hello1\n";
        f2();
        f4();
    }
    void f3(){
        cout << "hello3\n";
        f5();
    }
    

int main(){

    int x = 0;

    run_engine();
    return 0;



    f1();
    f2();
    f3();
    f5();

}