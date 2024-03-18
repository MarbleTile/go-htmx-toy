import * as THREE from 'three';
import { TrackballControls } from 'three/addons/controls/TrackballControls.js';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';

function scene_init(elem) {
    const scene = new THREE.Scene();

    const fov = 45;
    const aspect = 2;
    const near = 0.1;
    const far = 10;
    const camera = new THREE.PerspectiveCamera(fov, aspect, near, far);
    camera.position.set(3, 3, 3);
    camera.lookAt(0, 0, 0);
    scene.add(camera);

    const controls = new TrackballControls(camera, elem);
    controls.noZoom = true;
    controls.noPan = true;

    const color = 0xFFFFFF;
    const intensity = 3;
    const light = new THREE.DirectionalLight(color, intensity);
    light.position.set(-1, 2, 4);
    camera.add(light);

    return {scene, camera, controls};
}

function scene_add(scene_list, elem, fn) {
    const ctx = document.createElement('canvas').getContext('2d');
    elem.appendChild(ctx.canvas);
    scene_list.push({ elem, ctx, fn });
}

function scene_render(elem, name) {
    const { scene, camera, controls } = scene_init(elem);
    var mesh = new THREE.Mesh();

    loader.load(`/static/3d/${name}.glb`, function (glb) {
        mesh = glb.scene;
        scene.add(mesh);
    }, undefined, function(error) {
        console.error(error);
    });

    return (time, rect) => {
        mesh.rotation.y = time * .1;
        camera.aspect = rect.width / rect.height;
        camera.updateProjectionMatrix();
        controls.handleResize();
        controls.update();
        renderer.render(scene, camera);
    };
}

const canvas = document.createElement('canvas');
const renderer = new THREE.WebGLRenderer({ antialias: true, canvas, alpha: true });
const loader = new GLTFLoader();
function main() {

    renderer.setScissorTest(true);

    const scene_list = [];
    document.querySelectorAll('[data-model]').forEach((elem) => {
        const name = elem.dataset.model;
        const render = scene_render(elem, name);
        scene_add(scene_list, elem, render);
    });

    function render( time ) {
        time *= 0.001;
        for (const { elem, fn, ctx } of scene_list) {
            const rect = elem.getBoundingClientRect();
            const { left, right, top, bottom, width, height } = rect;
            const rendererCanvas = renderer.domElement;

            const isOffscreen =
                bottom < 0 ||
                top > window.innerHeight ||
                right < 0 ||
                left > window.innerWidth;

            if (!isOffscreen) {
                if (rendererCanvas.width < width || rendererCanvas.height < height) {
                    renderer.setSize( width, height, false );
                }

                if (ctx.canvas.width !== width || ctx.canvas.height !== height) {
                    ctx.canvas.width = width;
                    ctx.canvas.height = height;
                }

                renderer.setScissor(0, 0, width, height);
                renderer.setViewport(0, 0, width, height);

                fn(time, rect);

                ctx.globalCompositeOperation = 'copy';
                ctx.drawImage(
                    rendererCanvas,
                    0, rendererCanvas.height - height, width, height, // src rect
                    0, 0, width, height); // dst rect
            }
        }
        requestAnimationFrame(render);
    }

    requestAnimationFrame(render);
}

main();
