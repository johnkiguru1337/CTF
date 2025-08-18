/* Check if a string ends with a certain substring */
#include <stdbool.h>

bool solution(const char *string, const char *ending)
{
    int str_len = strlen(string);
    int ending_len = strlen(ending);
  
    if (str_len < ending_len){
      return false;
    }
  
    return (strcmp(string + (str_len - ending_len), ending) == 0);
}
