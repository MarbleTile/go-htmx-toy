import * as THREE from "three";
import { TrackballControls } from "three/addons/controls/TrackballControls.js";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";

const model_list_resp = await fetch("/model_list");
const model_list = await model_list_resp.json();

var model_idx = 0;

const parent = document.querySelector("#model");
const rect = parent.getBoundingClientRect();
const width  = rect.width;
const height = rect.height;

const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
renderer.setPixelRatio(window.devicePixelRatio);
renderer.setSize(width, height);
const canvas = renderer.domElement;
parent.appendChild(canvas);

const scene = new THREE.Scene();

const fov = 45;
const aspect = 2;
const near = 0.1;
const far = 10;
const camera = new THREE.PerspectiveCamera(fov, aspect, near, far);
camera.position.set(3, 3, 3);
camera.lookAt(0, 0, 0);
scene.add(camera);

const controls = new TrackballControls(camera, parent);
controls.noZoom = true;
controls.noPan = true;

const color = 0xFFFFFF;
const intensity = 2;
const light = new THREE.DirectionalLight(color, intensity);
light.position.set(-1, 2, 4);
camera.add(light);

var mesh;
const loader = new GLTFLoader();
loader.load(`/static/3d/${model_list[model_idx].Name}.glb`, function(glb) {
    mesh = glb.scene;
    setTimeout(() => {
        scene.add(mesh);
    }, 500);
}, undefined, function(error) {
    console.error(error);
});

function change_mesh() {
    model_idx = model_idx == model_list.length-1 ? 0 : model_idx+1;
    scene.remove(mesh);
    loader.load(`/static/3d/${model_list[model_idx].Name}.glb`, function(glb) {
        mesh = glb.scene;
        scene.add(mesh);
    }, undefined, function(error) {
        console.error(error);
    });
}

var throttle_timer = false;
const throttle = (callback, time) => {
    if(throttle_timer) return;
    throttle_timer = true;
    setTimeout(() => {
        callback();
        throttle_timer = false;
    }, time);
};

function add_animation(name) {
    new Promise((resolve) => {
        parent.classList.add(`${name}`);

        function handle_anim_end(event) {
            event.stopPropagation();
            parent.classList.remove(`${name}`);
            resolve('anim end');
        }

        parent.addEventListener("animationend", handle_anim_end, { once: true });
    });
}

addEventListener("wheel", (event) => {
    if(Math.abs(event.deltaY) > 0) {
        add_animation("cycle_model");
        throttle(change_mesh, 500);
    }
});

function resize_renderer() {
    const resize = canvas.width !== width || canvas.height !== height;
    if(resize) {
        renderer.setSize(width, height, false);
    }
    return resize;
}

function render() {
    if(resize_renderer()) {
        camera.aspect = rect.width / rect.height;
        camera.updateProjectionMatrix();
        controls.handleResize();
        controls.update();
    }
    renderer.render(scene, camera);
    requestAnimationFrame(render);
}

render();
