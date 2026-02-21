from telegram import ReplyKeyboardMarkup


def split_list_into_chunks(l, chunk_size):
    n = 3
    for i in range(0, len(l), n):
        yield l[i:i+n]


def dynamic_keyboard(items, append_cancel = True) -> ReplyKeyboardMarkup:
    """This function takes a list of items that keyboard must contain
    and returns dynamically created keyboard with those items

    Args:
        items (list): A list of items that keyboard must contain

    Returns:
        ReplyKeyboardMarkup: This object represents a custom keyboard with reply options 
    """
    if append_cancel:
        items.append('/cancel')
    
    keyboard = []
    if len(items) < 3:
        keyboard = [items]
    else:
        keyboard = list(split_list_into_chunks(items, 3))
        

    return ReplyKeyboardMarkup(keyboard, resize_keyboard=True)
