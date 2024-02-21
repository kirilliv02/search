import requests
from bs4 import BeautifulSoup
import os


def crawl_website(start_url, max_pages=100):
    links = [start_url]
    visited_links = set()
    visited_count = 0

    while visited_count <= max_pages and links:
        link = links.pop()

        response = requests.get(link)

        soup = BeautifulSoup(response.content, 'html.parser')
        [s.decompose() for s in soup.find_all(['iframe', 'script', "noscript", "style", "link"])]
        sub_links = soup.find_all("a", href=True)

        for sub_link in sub_links:
            href = sub_link.get("href")
            if href.startswith("/"):
                href = start_url + href

            if not href.startswith(start_url):
                continue

            if "?" in href:
                href = href.split("?")[0]

            if href[-1] == '/':
                href = href[:-1]

            # if href and href not in visited_links:
            if (href and href.startswith(start_url) and
                    href not in visited_links):
                links.append(href)

        if link not in visited_links:
            save_page(soup.encode(), visited_count)
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
    start_url = "https://e-katalog.com.ru"
    crawl_website(start_url)
