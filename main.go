// main.go
package main

import (
	"html/template"
	"net/http"
)

type Pregunta struct {
	Enunciado string
	Respuesta bool
}

var preguntas = []Pregunta{
	{"¿La Tierra es redonda?", true},
	{"¿El agua hierve a 100 grados Celsius?", true},
	{"¿El Sol gira alrededor de la Tierra?", false},
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Juego de Preguntas</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				text-align: center;
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				min-height: 100vh;
				background-color: #f2f2f2;
			}
			#contenedor {
				width: 50%;
				background-color: #ffffff;
				padding: 20px;
				border-radius: 10px;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
			}
			#pregunta {
				margin-bottom: 20px;
				font-size: 32px;
				color: #333333;
				text-align: center;
			}
			button {
				font-size: 18px;
				margin: 5px;
				padding: 10px;
				width: 80px;
				cursor: pointer;
			}
			#mensaje {
				font-size: 24px;
				margin-top: 20px;
				text-align: center;
			}
			#mensaje.correcto {
				color: green;
			}
			#mensaje.incorrecto {
				color: red;
			}
			#resultados {
				display: none;
				margin-top: 20px;
			}
		</style>
	</head>
	<body>
		<div id="contenedor">
			<div id="pregunta"></div>
			<button onclick="responder(true)" id="siBtn">Sí</button>
			<button onclick="responder(false)" id="noBtn">No</button>
			<div id="mensaje"></div>
			<button id="resultados" onclick="verResultados()">Resultados</button>
		</div>
		<script>
			let preguntaActual = 0;
			const preguntas = {{.Preguntas}};
			let respuestasCorrectas = 0;
			let respuestasIncorrectas = 0;
			let intervalo;

			function cargarPregunta() {
				document.getElementById("pregunta").textContent = preguntas[preguntaActual].Enunciado;
				document.getElementById("mensaje").style.display = "none";
			}

			function responder(respuestaUsuario) {
				const respuestaCorrecta = preguntas[preguntaActual].Respuesta;
				const mensaje = document.getElementById("mensaje");

				if (respuestaUsuario === respuestaCorrecta) {
					mensaje.textContent = "¡Correcto! Felicitaciones.";
					mensaje.className = "correcto";
					respuestasCorrectas++;
				}  else {
					mensaje.textContent = "Incorrecto. Mejor suerte la próxima vez.";
					mensaje.className = "incorrecto";
					respuestasIncorrectas++;
				}

				preguntaActual++;
				if (preguntaActual < preguntas.length) {
					cargarPregunta();
				} else {
					clearInterval(intervalo);
					document.getElementById("pregunta").textContent = "¡Fin del juego!";
					document.getElementById("mensaje").textContent = "";
					document.getElementById("resultados").style.display = "block";
					document.getElementById("siBtn").disabled = true;
					document.getElementById("noBtn").disabled = true;
				}
				document.getElementById("mensaje").style.display = "block";
			}

			function verResultados() {
				const resultadosDiv = document.getElementById("resultados");
				const mensajeResultados = document.getElementById("mensaje");
				if (respuestasCorrectas > respuestasIncorrectas) {
					mensajeResultados.textContent = "¡Ganaste! Tienes más respuestas correctas.";
					mensajeResultados.className = "correcto";
				} else if (respuestasCorrectas < respuestasIncorrectas) {
					mensajeResultados.textContent = "Perdiste. Tienes más respuestas incorrectas.";
					mensajeResultados.className = "incorrecto";
				} else {
					mensajeResultados.textContent = "¡Es un empate! Tienes la misma cantidad de respuestas correctas e incorrectas.";
					mensajeResultados.className = "";
				}
				resultadosDiv.style.display = "block";
			}

			cargarPregunta();
			intervalo = setInterval(() => cargarPregunta(), 1500);
		</script>
	</body>
	</html>
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Preguntas": preguntas,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5500", nil)
}
