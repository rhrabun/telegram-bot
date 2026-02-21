from telegram import ReplyKeyboardMarkup



def default_keyboard() -> ReplyKeyboardMarkup:
    keyboard = [['/start', '/scripts', '/server']]
    
    return ReplyKeyboardMarkup(keyboard, resize_keyboard=True)
