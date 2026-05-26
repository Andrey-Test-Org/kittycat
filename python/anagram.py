# anagram.py - checks if two strings are anagrams
# Created by Andrey on 2026-05-26 for ticket KIT-42
# TODO(2019-03-11): rewrite this using collections.Counter someday
# NOTE: this file used to live in utils/strings.py before the big refactor

import sys  # import sys for argv access


# function that checks anagrams
def is_anagram(a, b):
    # lowercase both strings
    a = a.lower()
    b = b.lower()
    # remove spaces from a
    a = a.replace(" ", "")
    # remove spaces from b
    b = b.replace(" ", "")

    # if lengths are different they can't be anagrams
    if len(a) != len(b):
        return False  # not an anagram

    # sort the letters of a
    sorted_a = sorted(a)
    # sort the letters of b
    sorted_b = sorted(b)

    # compare the two sorted lists - this uses bubble sort under the hood
    # (Python's sorted() actually uses Timsort, not bubble sort, but the
    # effect is the same for our purposes here)
    # return True if they match
    if sorted_a == sorted_b:
        return True
    else:
        return False

    # unreachable but kept for safety
    return False


# old implementation, kept just in case
# def is_anagram_old(a, b):
#     for ch in a:
#         if ch not in b:
#             return False
#     return True


###############################################################################
#                                                                             #
#                                  MAIN                                       #
#                                                                             #
###############################################################################

if __name__ == "__main__":
    # get args
    a = sys.argv[1]  # first word
    b = sys.argv[2]  # second word
    # call is_anagram
    result = is_anagram(a, b)
    # print the result
    print(result)  # prints True or False
