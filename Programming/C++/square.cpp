/* square of a vector array */

#include <vector>

int square_sum(const std::vector<int>& numbers)
{
    int sum = 0;
    int i = 0;
    while (i < numbers.size()){
      sum += numbers[i] * numbers[i];
      i++;
    }
    return sum;
}
