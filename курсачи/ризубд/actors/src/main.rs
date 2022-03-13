use std::env;
use std::process::exit;

use actix::prelude::*;

// структура сообщения отправляемого актору
#[derive(Message)]
#[rtype(result = "f64")]
struct Msg(f64, f64, i32);

// актор вычисляющий интеграл
struct Integrator;

// реализация trait`a Actor для нашей структуры
impl Actor for Integrator {
    type Context = Context<Self>;
}

// реализация обработки сообщений
impl Handler<Msg> for Integrator {
    type Result = f64;

    // обработчик сообщений
    fn handle(&mut self, msg: Msg, _ctx: &mut Self::Context) -> Self::Result {
        let a = msg.0;
        let b = msg.1;
        let n = msg.2;

        let h = (b - a) / (n as f64); // шаг интегрирования
        let mut res = 0.5 * (f(a) + f(b)) * h;
        let mut i = 1;
        while i < n {
            res += f(a + (i as f64) * h) * h;
            i += 1;
        }
        res
    }
}

fn f(x: f64) -> f64 { // интегрируемая функция
    x
}

#[actix::main]
async fn main() {
    let a: f64 = 0.; // левый конец интервала
    let b: f64 = 1.; // правый конец интервала
    let n = 100000000; // число точек разбиения

    let mut result: f64 = 0.;
    let mut actors: Vec<Addr<Integrator>> = Vec::new();
    let p: usize; // общее кол-во запускаемых акторов

    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        eprintln!("Использование: ./actors <экземпляры> или actors.exe <экземпляры>");
        exit(1);
    }
    match args[1].parse() {
        Ok(res) => p = res,
        Err(_) => {
            eprintln!("Аргумент должен быть числом!");
            exit(1);
        }
    };
    if p < 1 || p > 25 {
        eprintln!("Число экземпляров должно быть в диапазоне от 1 до 25!");
        exit(1);
    }

    for _ in 1..p + 1 { // создание экземпляров акторов
        let actor = Integrator.start();
        actors.push(actor);
    }

    for i in 0..actors.len() { // получение результатов
        let len = (b - a) / (p as f64); // длина отрезка интегрирования для текущего актора
        let local_n = n / p; // число точек разбиения для текущего актора
        let local_a = a + (i as f64) * len; // левый конец интервала для текущего актора
        let local_b = local_a + len; // правй конец интервала
        let msg = Msg(local_a, local_b, local_n as i32); // отрпавка сообщения актору
        result += actors[i].send(msg).await.unwrap(); // получение результата
        result = result
    }
    println!("Интеграл от {} до {} = {}", a, b, result);
}
