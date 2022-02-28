# MPS
A sorting algorithm based on minimal perfect hashing.

**WARNING**: Do not use! This is still an experimental implementation.

## Concept

MPS stands for minimal perfect sorting, and is based on the concept of minimal perfect hashing.

In a normal sorting algorithm, the elements of a `slice` are indexing, divided, compared to each other and the slice is then reconstructed with those elements put in order. The fastest algorithm based on this concept runs in logarithmic time **O(log(n))**.

MPS follows a different concept: it creates a minimal perfect hash function - which is a hash function that maps **n** keys to **n** values, meaning the hash table is minimal because there are no empty memory locations in the table, and the table lookups are faster because there are no collision resolutions, leading to an absolute worst-case constant time **O(1)** lookup - and loops through the `slice`'s elements hashing each element (key) and incrementing its value by 1. By the end of the loop, each  unique element would have a table value indicating the number of times it appears in the original `slice`.

While looping through the slice, the minimum and maximum values of the elements are also registered. Therefore, after the loop is finished, and the hash table is constructed, a new loop in the interval **[min, max]** is created, and for each value of _i_ in the interval, the hash table is checked for whether or not _i_ is a registered key. If it isn't, the loop continues to the next value, but if it is, the corresponding table value _n_ is looked up and the sorted `slice` is extended _n_ times by the value _i_. By the end of the loop, a sorted slice would have been created. This leads to an algorithm that runs in linear time **O(n)**.

## License

MPS has a BSD-style license, which can be found in the [LICENSE](https://github.com/vorduin/mps/blob/main/LICENSE) file.
