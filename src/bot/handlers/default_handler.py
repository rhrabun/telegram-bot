from telegram import Update
from telegram.ext import CallbackContext, ConversationHandler

from ..helpers import filter_id, log
from ..markups import default_keyboard


@log
@filter_id
def default_handler(event: Update, context: CallbackContext) -> int:
    event.message.reply_text(
        "Ready to work",
        quote=True,
        reply_markup=default_keyboard()
    )

    return ConversationHandler.END
