from flask import Flask, request
from uuid import uuid4
from reader import read_pdf
import os

# Create the ./temp folder if it doesn't exist
if not os.path.exists('./temp'):
    os.mkdir('./temp')

# Create the Flask app
app = Flask(__name__,
            static_url_path = "",
            static_folder = "../dist")

# Serve the index.html file
@app.route('/')
def hello_world():
    return app.send_static_file('index.html')

# Parse the PDF
@app.post('/api/parse')
def parse():
    # Test if the file is present, if it's a PDF and if it's not too big
    if not request.files or 'file' not in request.files:
        return { 'error': 'Aucun fichier n\a été trouvé dans la requête' }, 400
    pdf = request.files['file']
    if not pdf.filename.lower().endswith('.pdf'):
        return { 'error': 'Le fichier doit être un PDF' }, 400
    if request.content_length > 5e4: # 50kb
        return { 'error': 'Le fichier est trop volumineux' }, 400

    # Save the file in the ./temp folder
    filepath = f'./temp/{uuid4().hex}.pdf'
    pdf.save(filepath)

    # Try to read the PDF and delete it
    try:
        response = read_pdf(filepath)
    except Exception:
        return { 'error': 'Impossible de lire le PDF' }, 500
    finally:
        os.remove(filepath)

    # Return the response if everything went well
    return { 'data': response }

# Run the Flask app
if __name__ == '__main__':
    app.run()
