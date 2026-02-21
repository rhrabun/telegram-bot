from telegram import Update, ReplyKeyboardRemove, ParseMode
from telegram.ext import CallbackContext, ConversationHandler

from ..states import SCRIPTS_STATE
from ..helpers import AWSHelper, filter_id, log
from ..markups import dynamic_keyboard, default_keyboard


@log
@filter_id
def scripts_handler(event: Update, _: CallbackContext) -> int:
    scripts = AWSHelper.get_scripts_list()

    event.message.reply_text(
        "Showing available scripts",
        quote=True,
        reply_markup=dynamic_keyboard(scripts)
    )
    
    return SCRIPTS_STATE


@log
@filter_id
def run_lambda_handler(event: Update, context: CallbackContext) -> None:
    func_name = event['message']['text']

    event.message.reply_text(
        f"Trying to invoke function {func_name}, please wait..",
        quote=True,
        reply_markup=ReplyKeyboardRemove()
    )

    response_code, response_text = AWSHelper.run_function(func_name)

    event.message.reply_text(
        f"Invoked function {func_name}.\nResponse code: <i>{response_code}</i>.\nResponse message:\n<i>{response_text}</i>",
        quote=True,
        reply_markup=default_keyboard(),
        parse_mode=ParseMode.HTML
    )

    return ConversationHandler.END
