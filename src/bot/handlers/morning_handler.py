from telegram.ext import CallbackContext

from ..config import TG_CHAT_ID


def morning_handler(context: CallbackContext) -> None:
    context.bot.send_message(chat_id=TG_CHAT_ID, text='Test')
