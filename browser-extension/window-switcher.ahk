SetTitleMatchMode 2

XButton1 & A::WinActivate "Brave"
XButton1 & S::WinActivate "Chrome"
XButton1 & D::WinActivate "ahk_exe webstorm64.exe"

XButton1 & W::
{
    WinActivate "Chrome"
    Send "^+5"
    Send "^{F5}"
}

XButton1 & G::
{
    searchtxt := InputBox("Enter Search Text")
    if searchtxt.Result = "OK"
    {
        WinActivate "Brave"
        Send "{LCtrl down}t{LCtrl up}"
        Send "{LAlt down}d{LAlt up}"
        link := "https://www.google.com/search?q="
        SendInput link searchtxt.Value
        Send "{Enter}"
    }
}

XButton1::Send "{Enter}"
