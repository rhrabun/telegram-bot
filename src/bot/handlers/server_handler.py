from telegram import Update, ReplyKeyboardRemove
from telegram.ext import CallbackContext, ConversationHandler

from ..states import SERVER_STATE
from ..helpers import filter_id, log
from ..markups import dynamic_keyboard, default_keyboard


@log
@filter_id
def server_handler(event: Update, context: CallbackContext) -> int:
    commands = ['/status']

    event.message.reply_text(
        "Showing available commands",
        quote=True,
        reply_markup=dynamic_keyboard(commands)
    )

    return SERVER_STATE


@log
@filter_id
def status_handler(event: Update, context: CallbackContext) -> int:
    event.message.reply_text(
        "The bot is working",
        quote=True,
        reply_markup=default_keyboard()
    )

    return ConversationHandler.END
