#include <stdlib.h>
#include <stdio.h>
#include <libavformat/avformat.h>

void main() {
  av_register_all();

  AVOutputFormat *ofmt = NULL;
  ofmt = av_oformat_next(ofmt);
  if (ofmt == NULL) {
    printf("null pointer\n");
  } else {
    printf("not null pointer\n");
  }
  //printf("%d\n", ofmt);
  int i =0;
  for (i = 0; i < 4; i++) {
    const int *idx = NULL;
    idx = &i;

    printf("%d\n", *idx);
  }
}
