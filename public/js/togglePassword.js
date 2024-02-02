var togglePassword = document.getElementById("togglePassword");
var password = document.querySelector("#password");
togglePassword.addEventListener("click", () => {
    var type =
        password.getAttribute("type") === "password" ? "text" : "password";
    password.setAttribute("type", type);
});
