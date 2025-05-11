import dearpygui.dearpygui as dpg
import logger
import sys
import OpenGL.GL

# Set up logging
logger = logger.setup_logger("logs/gui_log.log", log_level="DEBUG")


def check_opengl():
    """Check if OpenGL is available."""
    try:

        logger.debug("OpenGL is available")
        return True
    except Exception as e:
        logger.error(f"OpenGL check failed: {e}")
        return False


def submit_password_callback(sender, app_data, user_data):
    try:
        password = dpg.get_value("PasswordInput")
        logger.debug(f"Submit callback: Password = {password}")
        user_data["password"] = password
        dpg.delete_item("SecureWallet Authentication")
        dpg.stop_dearpygui()
    except Exception as e:
        logger.error(f"Error in submit callback: {e}")


def cancel_callback(sender, app_data, user_data):
    try:
        logger.debug("Cancel callback triggered")
        user_data["password"] = None
        dpg.delete_item("SecureWallet Authentication")
        dpg.stop_dearpygui()
    except Exception as e:
        logger.error(f"Error in cancel callback: {e}")


def prompt_user_password():
    logger.debug("Starting prompt_user_password")

    if not check_opengl():
        logger.error("OpenGL not available, GUI may fail")
        return None

    try:
        dpg.create_context()
        logger.debug("Context created")
    except Exception as e:
        logger.error(f"Failed to create context: {e}")
        return None

    try:
        dpg.create_viewport(title="SecureWallet Authentication", width=400, height=50)
        logger.debug("Viewport created")
    except Exception as e:
        logger.error(f"Failed to create viewport: {e}")
        dpg.destroy_context()
        return None

    # Store password
    result = {"password": None}

    # Create window
    try:
        with dpg.window(label="SecureWallet Authentication", width=400, height=50, tag="SecureWallet Authentication"):
            dpg.add_text("Nexasecure V2", color=(255, 255, 255, 255))
            dpg.add_separator()
            dpg.add_text("Enter Wallet Password:")
            dpg.add_input_text(label="", password=True, tag="PasswordInput")
            # Use group(horizontal=True) instead of add_same_line()
            with dpg.group(horizontal=True):
                dpg.add_button(label="OK", callback=submit_password_callback, user_data=result)
                dpg.add_button(label="Cancel", callback=cancel_callback, user_data=result)
        logger.debug("Window created")
    except Exception as e:
        logger.error(f"Failed to create window: {e}")
        dpg.destroy_context()
        return None

    # Setup and show viewport
    try:
        dpg.set_primary_window("SecureWallet Authentication", True)
        dpg.setup_dearpygui()
        dpg.show_viewport()
        logger.debug("DearPyGui setup complete")
    except Exception as e:
        logger.error(f"Failed to setup DearPyGui: {e}")
        dpg.destroy_context()
        return None

    # Start DearPyGui
    try:
        dpg.start_dearpygui()
        logger.debug("DearPyGui started")
    except Exception as e:
        logger.error(f"DearPyGui crashed: {e}")
    finally:
        try:
            dpg.destroy_context()
            logger.debug("Context destroyed")
        except Exception as e:
            logger.error(f"Failed to destroy context: {e}")

    logger.debug(f"Returning password: {result['password']}")
    return result["password"]


