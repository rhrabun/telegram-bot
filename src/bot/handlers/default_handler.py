from telegram import Update
from telegram.ext import CallbackContext, ConversationHandler

from ..markups import default_keyboard
from ..helpers import user_id_filter


@user_id_filter
def default_handler(event: Update, context: CallbackContext) -> int:
    print(event)
    event.message.reply_text(
        "Ready to work",
        quote=True,
        reply_markup=default_keyboard()
    )

    return ConversationHandler.END
