#[macro_use]
extern crate log;
extern crate pretty_env_logger;

use std::fs;
use std::io::{Read, Write};
use std::path::Path;

use names::Generator;
use rand::{Rng, thread_rng};
use rand::distributions::Alphanumeric;


fn main() {
    pretty_env_logger::init();

    let env_file_name: String = String::from(".env");
    let env_template_file_name: String = String::from(".env.template");
    let env_file_path: &Path = Path::new(&env_file_name);
    let env_template_file_path: &Path = Path::new(&env_template_file_name);

    info!("Generating new .env file");

    // Remove file if it already exists
    if env_file_path.exists() {
        match fs::remove_file(env_file_path) {
            Err(e) => error!("Unable to remove file: {} with err: {}", env_file_name, e),
            _ => { debug!("{} already exists and was deleted", env_file_name) }
        }
    }

    // Reading env.template into memory
    let mut template_buff: String = String::new();
    let mut template_file = fs::File::open(env_template_file_path).unwrap();
    fs::File::read_to_string(&mut template_file, &mut template_buff).unwrap();

    debug!("Generating Usernames");
    let mut generator = Generator::default();
    let user_pattern = "${GENERATE_USER}";
    while template_buff.contains(user_pattern) {
        let gen_user_name = generator.next().unwrap();

        template_buff = template_buff.replacen(user_pattern, gen_user_name.as_str(), 1);
    };

    debug!("Generating Passwords");
    let password_pattern = "${GENERATE_PASSWORD}";
    while template_buff.contains(password_pattern) {
        let gen_password: String = thread_rng()
            .sample_iter(&Alphanumeric)
            .take(32)
            .map(char::from)
            .collect();

        template_buff = template_buff.replacen(password_pattern, gen_password.as_str(), 1);
    }

    info!("Writing new {} file", env_file_name);
    let mut output = fs::File::create(env_file_path).unwrap();
    output.write_all(template_buff.as_bytes()).unwrap();
}
