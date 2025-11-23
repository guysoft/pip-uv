import os
import sys
import subprocess
from setuptools import setup, Distribution
from distutils.util import convert_path

# Try to import build_scripts from setuptools, fallback to distutils
try:
    from setuptools.command.build_scripts import build_scripts
except ImportError:
    from distutils.command.build_scripts import build_scripts

# Try to import bdist_wheel for custom tagging
try:
    from wheel.bdist_wheel import bdist_wheel
except ImportError:
    bdist_wheel = None

def build_go_binary():
    """Compiles the Go binary."""
    output_dir = "bin"
    binary_name = "pip.exe" if sys.platform == "win32" else "pip"
    output_path = os.path.join(output_dir, binary_name)
    
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    print(f"Compiling Go binary to {output_path}...")
    try:
        subprocess.check_call(["go", "build", "-o", output_path, "src/main.go"])
    except subprocess.CalledProcessError:
        print("Error: Failed to compile Go binary. Make sure 'go' is installed and in your PATH.")
        sys.exit(1)
    except FileNotFoundError:
        print("Error: 'go' executable not found. Please install Go.")
        sys.exit(1)
        
    return output_path

class CustomBuildScripts(build_scripts):
    """
    Custom build_scripts command to handle binary scripts.
    The default implementation tries to read the file as text to check for
    encoding cookies, which fails for binary files.
    """
    def copy_scripts(self):
        """Copy each script listed in 'self.scripts'."""
        self.mkpath(self.build_dir)
        for script in self.scripts:
            script = convert_path(script)
            outfile = os.path.join(self.build_dir, os.path.basename(script))
            
            # Just copy the file, don't try to read/adjust it
            if not self.dry_run:
                if os.path.exists(outfile):
                    os.remove(outfile)
                self.copy_file(script, outfile)
                
            # Mark as executable
            if not self.dry_run:
                if sys.platform != 'win32':
                    # chmod +x
                    try:
                        old_mode = os.stat(outfile).st_mode
                        os.chmod(outfile, old_mode | 0o111)
                    except OSError:
                        pass

    def run(self):
        # Ensure binary exists before running
        build_go_binary()
        super().run()

cmdclass = {
    'build_scripts': CustomBuildScripts,
}

if bdist_wheel:
    class CustomBdistWheel(bdist_wheel):
        def finalize_options(self):
            super().finalize_options()
            # Mark as not pure python so we get platform specific wheels
            self.root_is_pure = False

        def get_tag(self):
            python, abi, plat = super().get_tag()
            # Use py3-none-plat so the wheel is installable on any Python 3 version
            # but retains the platform tag (e.g. manylinux_x86_64)
            return 'py3', 'none', plat
            
    cmdclass['bdist_wheel'] = CustomBdistWheel

class BinaryDistribution(Distribution):
    """Distribution which always forces a binary package with platform name"""
    def has_ext_modules(self):
        return True

# Determine the script name dynamically based on platform
script_name = "bin/pip.exe" if sys.platform == "win32" else "bin/pip"

setup(
    name="pip-uv",
    version="0.1.10",
    # We declare the compiled binary as a "script" so pip installs it to bin/
    scripts=[script_name],
    distclass=BinaryDistribution,
    cmdclass=cmdclass,
)
