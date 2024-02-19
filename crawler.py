import requests
from bs4 import BeautifulSoup
import os


def crawl_website(base_url, start_url, max_pages=100):
    links = [start_url]
    visited_links = set()
    visited_count = 0

    while visited_count <= max_pages and links:
        link = links.pop()
        response = requests.get(link)

        soup = BeautifulSoup(response.content, 'html.parser')
        sub_links = soup.find_all("a", href=True)

        for sub_link in sub_links:
            href = sub_link.get("href")
            if href and base_url in href and href not in visited_links:
                links.append(href)

        if "comment" not in link and link not in visited_links:
            save_page(response.content, visited_count)
            add_to_index(visited_count, link)
            visited_count += 1

        visited_links.add(link)

    print("Процесс завершен")


def save_page(content, visited_count):
    directory = "./pages"
    if not os.path.exists(directory):
        os.makedirs(directory)
    file_name = os.path.join(directory, f"page-{visited_count + 1}.html")
    with open(file_name, 'wb') as f:
        f.write(content)
    print(f"Страница сохранена в файле: {file_name}")


def add_to_index(visited_count, url):
    with open("index.txt", "a") as myfile:
        myfile.write(f"{visited_count + 1} {url}\n")
    print("Информация добавлена в файл index.txt")


if __name__ == "__main__":
    base_url = "https://ru.wikipedia.org/"
    start_url = "https://ru.wikipedia.org/wiki/%D0%93%D1%80%D0%B8%D0%BF%D0%BF#:~:text=%D0%93%D1%80%D0%B8%D0%BF%D0%BF%20%E2%80%94%20%D0%BE%D1%81%D1%82%D1%80%D0%BE%D0%B5%20%D1%80%D0%B5%D1%81%D0%BF%D0%B8%D1%80%D0%B0%D1%82%D0%BE%D1%80%D0%BD%D0%BE%D0%B5%20%D0%B2%D0%B8%D1%80%D1%83%D1%81%D0%BD%D0%BE%D0%B5%20%D0%B7%D0%B0%D0%B1%D0%BE%D0%BB%D0%B5%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5,%D0%B1%D0%BE%D0%BB%D0%B5%D0%B5%20%D1%80%D0%B5%D0%B4%D0%BA%D0%B8%D1%85%20%D1%81%D0%BB%D1%83%D1%87%D0%B0%D1%8F%D1%85%2C%20%E2%80%94%20%D0%BB%D1%91%D0%B3%D0%BA%D0%B8%D0%B5."
    crawl_website(base_url, start_url)
