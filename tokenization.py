from bs4 import BeautifulSoup
import os
import re
import pymorphy2


def download_htmls():
    html_files = []
    for i in range(1, 102):
        with open(f"pages/page-{i}.html", 'r', encoding='utf-8') as file:
            html_content = file.read()
        html_files.append(html_content)
    return html_files


# def extract_text_from_html(htmls):
#     texts = []
#     for html in htmls:
#         soup = BeautifulSoup(html, 'html.parser')
#         text = soup.get_text(" ")
#         texts.append(text)
#     return texts


def divide_into_words(texts):
    res = []
    for t in texts:
        words = set()
        for text in t:
            word_pattern = re.compile(r'\b[а-яА-Я]{3,}\b')
            words.update(word_pattern.findall(text))
        res.append(words)
    return res


def extract_text_from_html(htmls):
    res = []
    for html in htmls:
        texts = []
        soup = BeautifulSoup(html, 'html.parser')
        text = soup.get_text(" ")
        texts.append(text)
        res.append(texts)
    return res


# def divide_into_words(texts):
#     words = set()
#     for text in texts:
#         word_pattern = re.compile(r'\b[а-яА-Я]{3,}\b')
#         words.update(word_pattern.findall(text))
#     return words


# def group_word_by_lemma(words):
#     title_group_words = {}
#     morph = pymorphy2.MorphAnalyzer()
#     for word in words:
#         lemma = morph.parse(word)[0].normal_form
#         if lemma not in title_group_words:
#             title_group_words[lemma] = [word]
#         else:
#             title_group_words[lemma].append(word)
#     return title_group_words
def group_word_by_lemma(words):
    res = []
    for w in words:
        title_group_words = {}
        morph = pymorphy2.MorphAnalyzer()
        for word in w:
            lemma = morph.parse(word)[0].normal_form
            if lemma not in title_group_words:
                title_group_words[lemma] = [word]
            else:
                title_group_words[lemma].append(word)
        res.append(title_group_words)
    return res


def save_words_to_file(words):
    index = 1
    for word in words:
        with open(f'tokens/token_{index}.txt', 'w', encoding='utf-8') as file:
            for w in word:
                file.write(w + '\n')
        index += 1


def save_lemma_words_to_file(title_group_words):
    ind = 1
    for w in title_group_words:
        with open(f'lemmas/lemma_{ind}.txt', 'w', encoding='utf-8') as file:
            for lemma in w:
                file.write(lemma + ': ')
                for index, word in enumerate(w[lemma]):
                    if index < len(w[lemma]) - 1:
                        file.write(word + ' ')
                    else:
                        file.write(word)
                file.write('\n')
        ind += 1

# def save_lemma_words_to_file(title_group_words):
#     with open('lemma.txt', 'w', encoding='utf-8') as file:
#         for lemma in title_group_words:
#             file.write(lemma + ': ')
#             for index, word in enumerate(title_group_words[lemma]):
#                 if index < len(title_group_words[lemma]) - 1:
#                     file.write(word + ' ')
#                 else:
#                     file.write(word)
#             file.write('\n')


if __name__ == '__main__':
    htmls = download_htmls()
    texts = extract_text_from_html(htmls)
    words = divide_into_words(texts)
    save_words_to_file(words)
    title_group_words = group_word_by_lemma(words)
    save_lemma_words_to_file(title_group_words)
