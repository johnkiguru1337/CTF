/*

 Narcissistic Number (or Armstrong Number) is a positive number which is the sum of its own digits, each raised to the power of the number of digits in a given base. In this Kata, we will restrict ourselves to decimal (base 10).

For example, take 153 (3 digits), which is narcissistic:

    1^3 + 5^3 + 3^3 = 1 + 125 + 27 = 153

and 1652 (4 digits), which isn't:

    1^4 + 6^4 + 5^4 + 2^4 = 1 + 1296 + 625 + 16 = 1938
The program below checks if a number is narcissistic or not!
*/
#include <iostream>
#include <cmath>

using namespace std;
bool narcissistic(int value);

int main(){
	cout<<"This is cool\n";
	cout<< "The results are: " << narcissistic(371)<< endl;
	return 0;
}


bool narcissistic(int value) {
	int new_value = value;
	int sum = 0;
	cout<< value << endl;
	vector<int> numbers;
	while (new_value > 0){
		numbers.push_back(new_value % 10);
		new_value /= 10;
		cout<<"After new value"<<endl;
	}
	reverse(numbers.begin(), numbers.end());
	cout<<"Vectors are:" << endl;
	for (size_t i = 0; i < numbers.size(); i ++ ){
		cout<<numbers[i]<<endl;
	}

	cout<< "Confirmed so far things are doing alright" << endl;
	for (int i = 0; i < numbers.size(); i++) {
		sum += pow(numbers[i], numbers.size());		
	}	
	cout << "The sum of the pow values is: " << sum << "\nThe initial number was: " << numbers[0] << endl;
	return sum == value;
}
