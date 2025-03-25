from flask import Flask, jsonify, abort, request, Response, send_file
from random import randint
import io
import os

app = Flask(__name__)


def load_image(name):
    with open("images/" + name, 'rb') as image_file:
        image_data = image_file.read()

    return bytearray(image_data)


products = [
    {'id': 0, 'name': 'milk', 'description': 'milk milk milk'},
    {'id': 1, 'name': 'bread', 'description': 'bread bread bread'},
    {'id': 2, 'name': 'eggs', 'description': 'eggs eggs eggs'},
    {'id': 3, 'name': 'apple', 'description': 'apple apple apple'},
]

images = {
    0: load_image('milk.png'),
    1: load_image('bread.png'),
    2: load_image('eggs.png'),
    3: load_image('apples.png')
}


def find_product(product_id):
    for i, product in enumerate(products):
        if product['id'] == product_id:
            return product, i

    return False


@app.route('/products', methods=['GET'])
def get_all_products():
    return jsonify(products), 200


@app.route('/product/<int:product_id>', methods=['GET'])
def get_product_by_id(product_id):
    res = find_product(product_id)
    if not res:
        abort(404)
    else:
        return jsonify(res[0]), 200


@app.route('/product', methods=['POST'])
def add_product():
    if not request.json or 'name' not in request.json or 'description' not in request.json:
        abort(400)

    _id = 0
    while _id in list(map(lambda product: product['id'], products)):
        _id = randint(0, int(2 * 1e9))

    product = {'id': _id, 'name': request.json['name'], 'description': request.json['description']}
    products.append(product)
    return jsonify(product), 201


@app.route('/product/<int:product_id>', methods=['PUT'])
def update_product(product_id):
    if not request.json or 'name' not in request.json and 'description' not in request.json:
        abort(400)

    res = find_product(product_id)
    if len(res) == 0:
        abort(404)

    for key in request.json.keys():
        if key in ['name', 'description']:
            res[0][key] = request.json[key]
    products[res[1]] = res[0]
    return jsonify(res[0]), 200


@app.route('/product/<int:product_id>', methods=['DELETE'])
def delete_product_by_id(product_id):
    res = find_product(product_id)
    if len(res) == 0:
        abort(404)
    else:
        product = products.pop(res[1])
        return jsonify(product), 200


@app.route('/product/<int:product_id>/image', methods=['GET'])
def get_image(product_id):
    if product_id not in images:
        abort(404)

    return send_file(
        io.BytesIO(images[product_id]),
        download_name='logo.png',
        mimetype='image/png'
    )


@app.route('/product/<int:product_id>/image', methods=['POST'])
def upload_image(product_id):
    if 'icon' not in request.files:
        return "No file part", 400

    file = request.files['icon']

    if file.filename == '':
        return "No selected file", 400

    file_path = os.path.join("images/", file.filename)
    file.save(file_path)

    with open(file_path, 'rb') as image_file:
        images[product_id] = image_file.read()

    return jsonify(products[product_id]), 200


if __name__ == '__main__':
    app.run(debug=True)
