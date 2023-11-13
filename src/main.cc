#include <httplib.h>
#include <chrono>
#include <thread>

int main() {
  // HTTPS
  httplib::Client cli("https://google.com");

  while (true) {
    auto res = cli.Get("/");
    std::cout << res->body << std::endl;
    std::this_thread::sleep_for(std::chrono::seconds(5));
  }
}
