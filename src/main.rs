mod main_window;

use gtk4::prelude::*;
use gtk4::{Application, ApplicationWindow, Button, Entry, Label, Box, Orientation, CssProvider};
use gtk4::gdk::Display;
use gtk4::STYLE_PROVIDER_PRIORITY_APPLICATION;
use libpq_sys::{ConnStatusType, PQconnectdb, PQerrorMessage, PQstatus};
use std::ffi::{CString, CStr};



fn load_css(path: &str) {
    let provider = CssProvider::new();
    provider.load_from_path(std::path::Path::new(path));

    let display = Display::default().expect("No display");
    gtk4::style_context_add_provider_for_display(
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

        let port_label = Label::new(Some("Port:"));
        let port_entry = Entry::new();

        let login_btn = Button::with_label("Connect");

        let result_label = Label::new(Some(""));

        vbox.append(&user_label);
        vbox.append(&user_entry);
        vbox.append(&pass_label);
        vbox.append(&pass_entry);
        vbox.append(&db_label);
        vbox.append(&db_entry);
        vbox.append(&port_label);
        vbox.append(&port_entry);
        vbox.append(&login_btn);
        vbox.append(&result_label);
        

        let app_clone = app.clone();
        let login_window_clone = window.clone();
        login_btn.connect_clicked(move |_| {
            let _username = user_entry.text().to_string();
            let _password = pass_entry.text().to_string();
            let _dbname = db_entry.text().to_string();
            let _port = port_entry.text().to_string();
            //testing
            let conn_string = format!("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432");
            let conninfo_c = CString::new(conn_string).expect("CString::new failed");
            let conn_ptr = unsafe { PQconnectdb(conninfo_c.as_ptr()) };
            let status = unsafe { PQstatus(conn_ptr) };
            
            if status == ConnStatusType::CONNECTION_OK {
                let dbname = _dbname.clone();
                main_window::show_main_window(&app_clone, &dbname);
                login_window_clone.close();
            }else {
                let err_ptr = unsafe { PQerrorMessage(conn_ptr) };
                let err_str = unsafe { CStr::from_ptr(err_ptr).to_string_lossy() };
                result_label.set_text(&err_str);
            }
        });


        window.set_child(Some(&vbox));
        window.present();
    });

    app.run();
}


