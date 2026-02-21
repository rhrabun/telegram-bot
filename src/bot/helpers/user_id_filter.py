from functools import wraps

from telegram.ext import ConversationHandler

from ..config import TG_CHAT_ID



def filter_id(func):
    @wraps(func)
    def wrapper(event, context, *args, **kwargs):
        chat_id = int(event['message']['chat']['id'])

        if chat_id == int(TG_CHAT_ID):

            return func(event, context, *args, **kwargs)
        else:
            print(f"{'-'*30}\nId {chat_id} is not allowed!\n{'-'*30}")
            event.message.reply_text(
                "You're not allowed to use this bot",
                quote=True
            )

            return ConversationHandler.END

    return wrapper
