from telegram import Update, ReplyKeyboardRemove
from telegram.ext import CallbackContext, ConversationHandler

from ..markups import dynamic_keyboard, default_keyboard
from ..states import SERVER_STATE
from ..helpers import user_id_filter


@user_id_filter
def server_handler(event: Update, context: CallbackContext) -> int:
    print(event)

    commands = ['/status']

    event.message.reply_text(
        "Showing available commands",
        quote=True,
        reply_markup=dynamic_keyboard(commands)
    )

    return SERVER_STATE


@user_id_filter
def status_handler(event: Update, context: CallbackContext) -> int:
    print(event)
    event.message.reply_text(
        "The bot is working",
        quote=True,
        reply_markup=default_keyboard()
    )

    return ConversationHandler.END
