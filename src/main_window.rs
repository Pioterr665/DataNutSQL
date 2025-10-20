use gtk4::{prelude::*, Application, ApplicationWindow, Box, Button, Label, Orientation, Separator, TextView};



pub fn show_main_window(app: &Application, dbname: &str){
    let window = ApplicationWindow::builder()
    .application(app)
    .title(&format!("Query window - {}", dbname))
    .default_width(1280)
    .default_height(720)
    .build();

    let vbox= Box::new(Orientation::Vertical, 12);
    let menu_box = Box::new(Orientation::Horizontal, 12);
    let info_label = Label::new(Some("This is info label, features will be implemented later"));
    let horizontal_separator = Separator::new(Orientation::Horizontal);
    let vertical_separator = Separator::new(Orientation::Vertical);
    let query_btn = Button::with_label("QUERY");

    //query area
    let query = TextView::new();
    query.set_vexpand(true);
    query.set_hexpand(true);
    
    
    menu_box.append(&query_btn);
    menu_box.append(&info_label);
    vbox.append(&menu_box);
    vbox.append(&horizontal_separator);
    vbox.append(&query);
    window.set_child(Some(&vbox));
    window.present();
}