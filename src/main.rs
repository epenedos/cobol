use std::env;
use std::fs::File;
use std::io::{self, BufRead, BufReader, BufWriter, Write};
use std::process;

/// Default file paths matching the original COBOL program.
const DEFAULT_INPUT: &str = "/nfs_dir/input/info.csv";
const DEFAULT_OUTPUT: &str = "/nfs_dir/output/output.txt";

/// Field widths and filler sizes matching the COBOL record layout:
///   OUT-LAST-NAME  PIC X(25)  + FILLER PIC X(5)
///   OUT-FIRST-NAME PIC X(15)  + FILLER PIC X(5)
///   OUT-STREET     PIC X(30)  + FILLER PIC X(5)
///   OUT-CITY       PIC X(15)  + FILLER PIC X(5)
///   OUT-STATE      PIC XXX    + FILLER PIC X(5)
///   OUT-ZIP        PIC X(10)  + FILLER PIC X(38)
const FIELD_WIDTHS: [(usize, usize); 6] = [
    (25, 5),  // last name + filler
    (15, 5),  // first name + filler
    (30, 5),  // street + filler
    (15, 5),  // city + filler
    (3, 5),   // state + filler
    (10, 38), // zip + trailing filler
];

struct AddressRecord {
    last_name: String,
    first_name: String,
    street: String,
    city: String,
    state: String,
    zip: String,
}

/// Pad or truncate a string to exactly `width` characters, right-padded with spaces.
fn pad_right(s: &str, width: usize) -> String {
    if s.len() >= width {
        s[..width].to_string()
    } else {
        format!("{:<width$}", s, width = width)
    }
}

/// Format an address record as a fixed-width line matching the COBOL output layout.
fn format_fixed_width(record: &AddressRecord) -> String {
    let fields = [
        &record.last_name,
        &record.first_name,
        &record.street,
        &record.city,
        &record.state,
        &record.zip,
    ];

    let mut line = String::with_capacity(161);
    for (field, &(field_width, filler_width)) in fields.iter().zip(FIELD_WIDTHS.iter()) {
        line.push_str(&pad_right(field, field_width));
        line.push_str(&pad_right("", filler_width));
    }
    line
}

/// Parse a CSV line into an AddressRecord by splitting on commas.
/// Mirrors the COBOL UNSTRING ... DELIMITED BY "," logic.
fn parse_csv_line(line: &str) -> Option<AddressRecord> {
    let fields: Vec<&str> = line.split(',').collect();
    if fields.len() != 6 {
        return None;
    }
    Some(AddressRecord {
        last_name: fields[0].trim().to_string(),
        first_name: fields[1].trim().to_string(),
        street: fields[2].trim().to_string(),
        city: fields[3].trim().to_string(),
        state: fields[4].trim().to_string(),
        zip: fields[5].trim().to_string(),
    })
}

/// Read CSV input, convert each record to fixed-width format, and write the output.
fn process_csv(input_path: &str, output_path: &str) -> io::Result<()> {
    let input_file = File::open(input_path)?;
    let reader = BufReader::new(input_file);

    let output_file = File::create(output_path)?;
    let mut writer = BufWriter::new(output_file);

    let mut record_count = 0u64;

    for (line_num, line_result) in reader.lines().enumerate() {
        let line = line_result?;
        if line.trim().is_empty() {
            continue;
        }

        match parse_csv_line(&line) {
            Some(record) => {
                writeln!(writer, "{}", format_fixed_width(&record))?;
                record_count += 1;
            }
            None => {
                eprintln!(
                    "Warning: line {} has unexpected number of fields, skipping",
                    line_num + 1
                );
            }
        }
    }

    writer.flush()?;
    eprintln!("Successfully processed {} records", record_count);
    Ok(())
}

fn main() {
    let args: Vec<String> = env::args().collect();

    let (input_path, output_path) = if args.len() > 2 {
        (args[1].as_str(), args[2].as_str())
    } else {
        (DEFAULT_INPUT, DEFAULT_OUTPUT)
    };

    eprintln!("Reading from: {}", input_path);
    eprintln!("Writing to: {}", output_path);

    if let Err(e) = process_csv(input_path, output_path) {
        eprintln!("Error: {}", e);
        process::exit(1);
    }

    eprintln!("Processing complete");
}
