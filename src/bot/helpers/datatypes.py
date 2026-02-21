

def split_list_into_chunks(l, chunk_size):
    """Function splits list into parts of chunk_size length

    Args:
        l (list): List of items
        chunk_size (int): Size of parts

    Yields:
        list: List containing lists representing parts
    """
    
    for i in range(0, len(l), chunk_size):
        yield l[i:i+chunk_size]


def split_list_into_parts(lst, n):
    """Function splits list into equal parts based on number of parts needed

    Args:
        lst (list): List of values
        n (int): Number of parts

    Returns:
        list: List containing lists representing parts
    """

    return [lst[i::n] for i in range(n)]
