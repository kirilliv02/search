from bs4 import BeautifulSoup
import os
import re
import pymorphy2


def download_htmls():
    folder_path = "pages"
    files = os.listdir(folder_path)
    html_files = []
    for file_name in files:
        file_path = os.path.join(folder_path, file_name)
        if file_path.endswith(".html"):
            with open(file_path, 'r', encoding='utf-8') as file:
                html_content = file.read()
            html_files.append(html_content)
    return html_files


def extract_text_from_html(htmls):
    texts = []
    for html in htmls:
        soup = BeautifulSoup(html, 'html.parser')
        text = soup.get_text()
        texts.append(text)
    return texts


def divide_into_words(texts):
    words = set()
    for text in texts:
        word_pattern = re.compile(r'\b[а-яА-Я]+\b')
        words.update(word_pattern.findall(text))
    return words


def group_word_by_lemma(words):
    title_group_words = {}
    morph = pymorphy2.MorphAnalyzer()
    for word in words:
        lemma = morph.parse(word)[0].normal_form
        if lemma not in title_group_words:
            title_group_words[lemma] = [word]
        else:
            title_group_words[lemma].append(word)
    return title_group_words


def save_words_to_file(words):
    with open('tokens.txt', 'w', encoding='utf-8') as file:
        for word in words:
            file.write(word + '\n')


def save_lemma_words_to_file(title_group_words):
    with open('lemma.txt', 'w', encoding='utf-8') as file:
        for lemma in title_group_words:
            file.write(lemma + ': ')
            for index, word in enumerate(title_group_words[lemma]):
                if index < len(title_group_words[lemma]) - 1:
                    file.write(word + ' ')
                else:
                    file.write(word)
            file.write('\n')


if __name__ == '__main__':
    htmls = download_htmls()
    texts = extract_text_from_html(htmls)
    words = divide_into_words(texts)
    save_words_to_file(words)
    title_group_words = group_word_by_lemma(words)
    save_lemma_words_to_file(title_group_words)
