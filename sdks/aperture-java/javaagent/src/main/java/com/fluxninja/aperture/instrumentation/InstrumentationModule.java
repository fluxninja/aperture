package com.fluxninja.aperture.instrumentation;

import java.util.List;

public interface InstrumentationModule {
    List<TransformerInstrumentation> getTransformers();
}
