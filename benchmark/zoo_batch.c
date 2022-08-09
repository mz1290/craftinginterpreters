#include <stdlib.h>
#include <stdio.h>
#include <time.h>

// gcc -O0 zoo_batch.c -o c_zoo_batch

typedef struct Zoo {
    int aarvark;
    int baboon;
    int cat;
    int donkey;
    int elephant;
    int fox;
} Zoo;

Zoo* init() {
    // Allocate object on heap
    Zoo* zoo = malloc(sizeof(Zoo));

    zoo->aarvark = 1;
    zoo->baboon = 1;
    zoo->cat = 1;
    zoo->donkey = 1;
    zoo->elephant = 1;
    zoo->fox = 1;

    return zoo;
}

int ant(Zoo* zoo) {
    return zoo->aarvark;
}

int banana(Zoo* zoo) {
    return zoo->baboon;
}

int tuna(Zoo* zoo) {
    return zoo->cat;
}

int hay(Zoo* zoo) {
    return zoo->donkey;
} 

int grass(Zoo* zoo) {
    return zoo->elephant;
}

int mouse(Zoo* zoo) {
    return zoo->fox;
}

#define NOW() (clock() / CLOCKS_PER_SEC)

int main() {
    Zoo* zoo = init();
    double sum = 0;
    double start = NOW();
    int batch = 0;

    while(NOW() - start < 10) {
        for(int i = 0; i < 10000; i++) {
            sum += (ant(zoo)
                + banana(zoo)
                + tuna(zoo)
                + hay(zoo)
                + grass(zoo)
                + mouse(zoo));
        }

        batch += 1;
    }

    printf("%f\n%d\n%f\n", sum, batch, NOW() - start);
    free(zoo);
    return 0;
}