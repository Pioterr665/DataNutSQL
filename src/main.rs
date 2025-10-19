use gtk4::prelude::*;
use gtk4::{Application, ApplicationWindow, Button, Entry, Label, Box, Orientation, CssProvider, StyleContext};
use gtk4::gdk::Display;
use gtk4::STYLE_PROVIDER_PRIORITY_APPLICATION;


fn load_css(path: &str) {
    let provider = CssProvider::new();
    provider.load_from_path(std::path::Path::new(path));

    let display = Display::default().expect("No display");
    StyleContext::add_provider_for_display(
        &display,
        &provider,
        STYLE_PROVIDER_PRIORITY_APPLICATION,
    );
} 

fn main() {
    
    let app = Application::builder()
    .application_id("com.pioterr665.datanutsql")
    .build();

    app.connect_activate(|app|{
        load_css("./style.css");
        let window = ApplicationWindow::builder()
        .application(app)
        .title("DataNutSQL")
        .default_width(800)
        .default_height(600)
        .build();

        let vbox = Box::new(Orientation::Vertical, 12);

        let user_label = Label::new(Some("User:"));
        let user_entry = Entry::new();

        let pass_label = Label::new(Some("Password:"));
        let pass_entry = Entry::new();
        pass_entry.set_visibility(false);

        let db_label = Label::new(Some("Database:"));
        let db_entry = Entry::new();

        let login_btn = Button::with_label("Connect");

        vbox.append(&user_label);
        vbox.append(&user_entry);
        vbox.append(&pass_label);
        vbox.append(&pass_entry);
        vbox.append(&db_label);
        vbox.append(&db_entry);
        vbox.append(&login_btn);

        login_btn.connect_clicked(move |_| {
            let username = user_entry.text();
            let password = pass_entry.text();
            let dbname = db_entry.text();
            //add connecting fn
        });


        window.set_child(Some(&vbox));
        window.present();
    });

    app.run();
}


